package nyaa

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
)

const nyaaConfig = "nyaa-config"

type config struct {
	NyaaID int64 `json:"NyaaID" bson:"NyaaID"`
	ByID   bool  `json:"ByID" bson:"ByID"`
	BySeed bool  `json:"BySeed" bson:"BySeed"`
}

type record struct {
	ID           string `json:"_id" bson:"_id"`
	ASW          config `json:"ASW" bson:"ASW"`
	BeatriceRaws config `json:"BeatriceRaws" bson:"BeatriceRaws"`
	Ember        config `json:"Ember" bson:"Ember"`
	EraiRaws     config `json:"EraiRaws" bson:"EraiRaws"`
	IrizaRaws    config `json:"IrizaRaws" bson:"IrizaRaws"`
	Judas        config `json:"Judas" bson:"Judas"`
	LostYears    config `json:"LostYears" bson:"LostYears"`
	Moozzi2      config `json:"Moozzi2" bson:"Moozzi2"`
	OhysRaws     config `json:"OhysRaws" bson:"OhysRaws"`
	ReinForce    config `json:"ReinForce" bson:"ReinForce"`
	SubsPlease   config `json:"SubsPlease" bson:"SubsPlease"`
	Yameii       config `json:"Yameii" bson:"Yameii"`
}

func (s *Nyaa) Init(ctx context.Context) error {
	var names = make([]string, 0)

	names = append(names, nyaaConfig)
	names = append(names, asw_toorent, asw_anime)
	names = append(names, beatrice_raws_toorent, beatrice_raws_anime)
	names = append(names, ember_toorent, ember_anime)
	names = append(names, erai_raws_toorent, erai_raws_anime)
	names = append(names, iriza_raws_toorent, iriza_raws_anime)
	names = append(names, judas_toorent, judas_anime)
	names = append(names, lost_years_toorent, lost_years_anime)
	names = append(names, moozzi2_toorent, moozzi2_anime)
	names = append(names, ohys_raws_toorent, ohys_raws_anime)
	names = append(names, rein_force_toorent, rein_force_anime)
	names = append(names, subs_please_toorent, subs_please_anime)
	names = append(names, yameii_toorent, yameii_anime)

	filter := bson.M{
		"name": bson.M{
			"$in": names,
		},
	}

	cursor, err := s.db.ListCollections(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	// Create a map to track which collections exist
	exists := make(map[string]bool, len(names))
	for _, name := range names {
		exists[name] = false
	}

	// Iterate through the matching collections and update the map
	for cursor.Next(ctx) {
		var collInfo struct {
			Name string `bson:"name"`
		}
		if err := cursor.Decode(&collInfo); err != nil {
			return err
		}
		exists[collInfo.Name] = true
	}

	for name, here := range exists {
		if !here {
			s.log.Warn("collection does not exist, creating it...", "name", name)
			if err := s.db.CreateCollection(ctx, name); err != nil {
				return err
			}

			indexModel := mongo.IndexModel{
				Keys: bson.D{
					bson.E{Key: "NyaaID", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			}

			_, err = s.db.Collection(name).Indexes().CreateOne(context.Background(), indexModel)
			if err != nil {
				s.log.Error("cannot create unique index", "error", err)
				return err
			}
			s.log.Warn("creating unique index ...")

			indexModel = mongo.IndexModel{
				Keys: bson.D{
					bson.E{Key: "Name", Value: "text"},
					bson.E{Key: "CRC", Value: "text"},
				},
				Options: options.Index(),
			}

			_, err = s.db.Collection(name).Indexes().CreateOne(context.Background(), indexModel)
			if err != nil {
				s.log.Error("cannot create search index", "error", err)
				return err
			}
			s.log.Warn("creating search index ...")

			exists[name] = true
		}
	}

	return nil
}

func (s *Nyaa) Save(ctx context.Context) {
	if s.record.ID == "" {
		return
	}
	s.log.Warn("saving the latest config update ...")

	filter := bson.D{{Key: "_id", Value: nyaaConfig}}
	update := bson.D{{Key: "$set", Value: s.record}}

	result := s.db.Collection(nyaaConfig).FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetUpsert(true))
	if err := result.Err(); err != nil {
		s.log.Error("cannnot update the config data", "error", err)
	}

	s.log.Info("config was updated successfuly")
}

func (s *Nyaa) load(ctx context.Context) error {
	filter := bson.M{"_id": nyaaConfig}

	err := s.db.Collection(nyaaConfig).FindOne(ctx, filter).Decode(s.record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.record.ID = nyaaConfig
			if _, err = s.db.Collection(nyaaConfig).InsertOne(ctx, s.record); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (s *Nyaa) Start(ctx context.Context) error {
	err := s.load(ctx)
	if err != nil {
		return err
	}

	group, ctx := errgroup.WithContext(ctx)
	s.run(ctx, group)

	tick := time.NewTicker(90 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			tick.Stop()
			return group.Wait()
		case t := <-tick.C:
			s.run(ctx, group)
			s.log.Info("last torrents update at", "time", t)
		}
	}
}

func (s *Nyaa) run(ctx context.Context, group *errgroup.Group) {
	time.Sleep(900 * time.Millisecond)
	group.Go(func() error {
		return s.scrape(ctx, s.db.Collection(asw_toorent), &s.record.ASW, asw_path)
	})

	time.Sleep(900 * time.Millisecond)
	group.Go(func() error {
		return s.scrape(ctx, s.db.Collection(beatrice_raws_toorent), &s.record.BeatriceRaws, beatrice_raws_path)
	})

	time.Sleep(900 * time.Millisecond)
	group.Go(func() error {
		return s.scrape(ctx, s.db.Collection(ember_toorent), &s.record.Ember, ember_path)
	})

	time.Sleep(900 * time.Millisecond)
	group.Go(func() error {
		return s.scrape(ctx, s.db.Collection(erai_raws_toorent), &s.record.EraiRaws, erai_raws_path)
	})

	time.Sleep(900 * time.Millisecond)
	group.Go(func() error {
		return s.scrape(ctx, s.db.Collection(iriza_raws_toorent), &s.record.IrizaRaws, iriza_raws_path)
	})

	time.Sleep(900 * time.Millisecond)
	group.Go(func() error {
		return s.scrape(ctx, s.db.Collection(judas_toorent), &s.record.Judas, judas_path)
	})

	time.Sleep(900 * time.Millisecond)
	group.Go(func() error {
		return s.scrape(ctx, s.db.Collection(lost_years_toorent), &s.record.LostYears, lost_years_path)
	})

	time.Sleep(900 * time.Millisecond)
	group.Go(func() error {
		return s.scrape(ctx, s.db.Collection(moozzi2_toorent), &s.record.Moozzi2, moozzi2_path)
	})

	time.Sleep(900 * time.Millisecond)
	group.Go(func() error {
		return s.scrape(ctx, s.db.Collection(ohys_raws_toorent), &s.record.OhysRaws, ohys_raws_path)
	})

	time.Sleep(900 * time.Millisecond)
	group.Go(func() error {
		return s.scrape(ctx, s.db.Collection(rein_force_toorent), &s.record.ReinForce, rein_force_path)
	})

	time.Sleep(900 * time.Millisecond)
	group.Go(func() error {
		return s.scrape(ctx, s.db.Collection(subs_please_toorent), &s.record.SubsPlease, subs_please_path)
	})

	time.Sleep(900 * time.Millisecond)
	group.Go(func() error {
		return s.scrape(ctx, s.db.Collection(yameii_toorent), &s.record.Yameii, yameii_path)
	})
}
