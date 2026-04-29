//go:build !nocry

package flags

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"time"

	cmdutils "github.com/digitalghost-dev/poke-cli/cmd/utils"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/ebitengine/oto/v3"
	"github.com/jfreymuth/oggvorbis"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func CryFlag(endpoint, pokemonName string) error {
	pokemonStruct, _, err := connections.PokemonApiCall(endpoint, pokemonName, connections.APIURL)
	if err != nil {
		return err
	}

	cryURL := pokemonStruct.Cries.Latest
	if cryURL == "" {
		return fmt.Errorf("%s", cmdutils.FormatError("No cry available for"+pokemonName))
	}

	fmt.Printf("Playing %s's cry...\n", cases.Title(language.English).String(pokemonName))

	resp, err := http.Get(cryURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	oggBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	dec, err := oggvorbis.NewReader(bytes.NewReader(oggBytes))
	if err != nil {
		return err
	}

	pcm := make([]float32, dec.Length()*int64(dec.Channels()))
	read := 0
	for read < len(pcm) {
		n, err := dec.Read(pcm[read:])
		read += n
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	pcm = pcm[:read]

	pcmBytes := new(bytes.Buffer)
	if err := binary.Write(pcmBytes, binary.LittleEndian, pcm); err != nil {
		return err
	}

	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   dec.SampleRate(),
		ChannelCount: dec.Channels(),
		Format:       oto.FormatFloat32LE,
	})
	if err != nil {
		return fmt.Errorf("%s", cmdutils.FormatError(fmt.Sprintf("Could not initialize audio: %v", err)))
	}
	<-ready

	p := ctx.NewPlayer(pcmBytes)
	p.Play()
	for p.IsPlaying() {
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(300 * time.Millisecond)

	return nil
}
