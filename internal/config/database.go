package config

import (
	"fmt"
	"log"
	"testProject/internal/model"

	"github.com/restream/reindexer/v4"
	_ "github.com/restream/reindexer/v4/bindings/cproto"
)

type NativeDatabase struct {
	DB *reindexer.Reindexer
}

var modelsMap = map[string]interface{}{
	"documents": model.Document{},
}

func Connect(conf *Configuration) (*NativeDatabase, error) {

	db := reindexer.NewReindex(conf.Database.GetConnectionString(), reindexer.WithCreateDBIfMissing())
	err := ensureNamespaces(db, modelsMap)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected!")

	return &NativeDatabase{db}, nil
}

func (db *NativeDatabase) Close() {
	db.DB.Close()
	println("Closing database...")
}

func ensureNamespaces(db *reindexer.Reindexer, namespaces map[string]interface{}) error {
	for name, m := range namespaces {
		if err := db.OpenNamespace(name, reindexer.DefaultNamespaceOptions(), m); err != nil {
			return fmt.Errorf("cannot open namespace %s: %w", name, err)
		}
	}
	return nil
}
