// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ctxgroup"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// processPortConnectionCSV imports port connection data (from CSV file to DB)
func (m *importer) processPortConnectionCSV(w http.ResponseWriter, r *http.Request) {
	log := m.log.For(r.Context())
	log.Debug("PortConnection- started")
	if err := r.ParseMultipartForm(maxFormSize); err != nil {
		log.Warn("parsing multipart form", zap.Error(err))
		http.Error(w, "cannot parse form", http.StatusUnprocessableEntity)
		return
	}

	ctx, mr := r.Context(), m.r.Mutation()
	for fileName := range r.MultipartForm.File {
		firstLine, reader, err := m.newReader(fileName, r)
		if err != nil {
			log.Warn("creating csv reader", zap.Error(err), zap.String("filename", fileName))
			http.Error(w, fmt.Sprintf("cannot handle file %q", fileName), http.StatusUnprocessableEntity)
			return
		}
		equipNameAIndex := findIndexForSimilar(firstLine, "A_Equipment")
		if equipNameAIndex == -1 {
			errorReturn(w, "Couldn't find 'A_Equipment' title", log, nil)
			return
		}

		equipNameBIndex := findIndexForSimilar(firstLine, "B_Equipment")
		if equipNameBIndex == -1 {
			errorReturn(w, "Couldn't find 'B_Equipment' title", log, nil)
			return
		}

		portAIndex := findIndexForSimilar(firstLine, "A_Port")
		if portAIndex == -1 {
			errorReturn(w, "Couldn't find 'A_Port' title", log, nil)
			return
		}

		portBIndex := findIndexForSimilar(firstLine, "B_Port")
		if portBIndex == -1 {
			errorReturn(w, "Couldn't find 'B_Port' title", log, nil)
			return
		}

		for {
			line, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Warn("cannot read row", zap.Error(err))
				continue
			}

			names := []string{
				line[portAIndex],
				line[portBIndex],
			}
			if names[0] == "" || names[1] == "" {
				continue
			}
			enames := []string{
				line[equipNameAIndex],
				line[equipNameBIndex],
			}

			ids := make([]string, len(names))
			{
				g := ctxgroup.WithContext(ctx)
				for i := range names {
					i := i
					g.Go(func(ctx context.Context) error {
						id, err := m.ClientFrom(ctx).EquipmentType.Query().
							QueryEquipment().
							Where(equipment.NameIn(
								enames[i],
								// Special case for cross-Layer's imports
								enames[i]+"CHCGILWK",
							)).
							QueryPorts().
							Where(equipmentport.HasDefinitionWith(
								equipmentportdefinition.NameIn(
									names[i],
									strings.ToLower(names[i]),
									strings.ToUpper(names[i]),
								),
							)).
							OnlyID(ctx)
						if err != nil {
							return errors.WithMessagef(err, "fetching one port: %q, equipment=%q", names[i], enames[i])
						}
						ids[i] = id
						return nil
					})
				}
				if err := g.Wait(); err != nil {
					errorReturn(w, err.Error(), log, err)
					continue
				}
			}

			pA := m.ClientFrom(ctx).EquipmentType.Query().
				QueryEquipment().
				QueryPorts().
				Where(equipmentport.ID(ids[0])).
				OnlyX(ctx)
			pB := m.ClientFrom(ctx).EquipmentType.Query().
				QueryEquipment().
				QueryPorts().
				Where(equipmentport.ID(ids[1])).
				OnlyX(ctx)

			if _, err = mr.AddLink(ctx,
				models.AddLinkInput{
					Sides: []*models.LinkSide{
						{Equipment: pA.QueryParent().OnlyXID(ctx), Port: pA.QueryDefinition().OnlyXID(ctx)},
						{Equipment: pB.QueryParent().OnlyXID(ctx), Port: pB.QueryDefinition().OnlyXID(ctx)},
					},
				}); err != nil {
				log.Warn("cannot connect ports",
					zap.Strings("ids", ids),
					zap.Strings("names", names),
					zap.Error(err),
				)
				continue
			}
			log.Info("connected ports",
				zap.Strings("ids", ids),
				zap.Strings("names", names),
			)
		}
	}
	log.Debug("PortConnection- Done")
}
