package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/drip"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared/components"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared/scenes"
)

func SerializeCallback(scene drip.Scene) ([]byte, error) {
	query := blueprint.Queries.ActionBuffer
	cursor := warehouse.Factory.NewCursor(query, scene.Storage())

	sEntities := []warehouse.SerializedEntity{}

	for range cursor.Next() {

		e, err := cursor.CurrentEntity()
		if err != nil {
			return nil, err
		}

		if !e.Valid() {
			log.Println("skipping invalid", e.Valid())
			continue
		}

		se := e.SerializeExclude(
			client.Components.SpriteBundle,
			client.Components.SoundBundle,
		)

		sEntities = append(sEntities, se)
	}

	serSto := warehouse.SerializedStorage{
		Entities:    sEntities,
		CurrentTick: scene.CurrentTick(),
		Version:     "net",
	}
	stateForJson, err := warehouse.PrepareForJSONMarshal(serSto)
	if err != nil {
		return nil, err
	}
	return json.Marshal(stateForJson)
}

func NewConnectionEntityCreate(conn drip.Connection, s drip.Server) (warehouse.Entity, error) {
	serverActiveScenes := s.ActiveScenes()

	if len(serverActiveScenes) == 0 {
		return nil, errors.New("No active scenes to find player in")
	}

	scene := serverActiveScenes[0]
	sto := scene.Storage()

	query := warehouse.Factory.NewQuery().And(components.PlayerSpawnComponent)
	cursor := warehouse.Factory.NewCursor(query, sto)

	var spawn components.PlayerSpawn

	for range cursor.Next() {
		match := components.PlayerSpawnComponent.GetFromCursor(cursor)
		spawn = *match
		break
	}

	return scenes.NewPlayer(spawn.X, spawn.Y, sto)
}
