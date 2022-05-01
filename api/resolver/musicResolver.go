package resolver

import (
	"github.com/danilomarques1/graphmusic/model"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

type MusicResolver struct {
	musicRepository model.MusicRepository
}

func NewMusicResolver(musicRepository model.MusicRepository) *MusicResolver {
	return &MusicResolver{musicRepository: musicRepository}
}

func (mr *MusicResolver) FindAll(p graphql.ResolveParams) (interface{}, error) {
	musics, err := mr.musicRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return musics, nil
}

func (mr *MusicResolver) FindByName(p graphql.ResolveParams) (interface{}, error) {
	name := p.Args["name"].(string)
	music, err := mr.musicRepository.FindByName(name)
	if err != nil {
		return nil, err
	}
	return music, nil
}

func (mr *MusicResolver) Save(p graphql.ResolveParams) (interface{}, error) {
	author := p.Args["author"].(string)
	name := p.Args["name"].(string)

	id := uuid.NewString()
	music := model.Music{Id: id, Author: author, Name: name}
	err := mr.musicRepository.Save(&music)
	if err != nil {
		return nil, err
	}
	return music, nil
}

func (mr *MusicResolver) RemoveByName(p graphql.ResolveParams) (interface{}, error) {
	name := p.Args["name"].(string)
	music, err := mr.musicRepository.RemoveByName(name)
	if err != nil {
		return nil, err
	}
	return music, nil
}

func (mr *MusicResolver) UpdateByName(p graphql.ResolveParams) (interface{}, error) {
	name := p.Args["name"].(string)
	musicMap := p.Args["music"].(map[string]interface{})
	nAuthor := musicMap["author"].(string)
	nName := musicMap["name"].(string)
	music := model.Music{Name: nName, Author: nAuthor}
	updated, err := mr.musicRepository.UpdateByName(name, &music)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (mr *MusicResolver) RemoveById(p graphql.ResolveParams) (interface{}, error) {
	id := p.Args["id"].(string)
	music, err := mr.musicRepository.RemoveById(id)
	if err != nil {
		return nil, err
	}
	return music, nil
}
