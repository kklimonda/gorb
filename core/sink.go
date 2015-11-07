/*
   Copyright (c) 2015 Andrey Sibiryov <me@kobology.ru>
   Copyright (c) 2015 Other contributors as noted in the AUTHORS file.

   This file is part of GORB - Go Routing and Balancing.

   GORB is free software; you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation; either version 3 of the License, or
   (at your option) any later version.

   GORB is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
   GNU Lesser General Public License for more details.

   You should have received a copy of the GNU Lesser General Public License
   along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package core

import (
	"github.com/kobolog/gorb/pulse"

	log "github.com/sirupsen/logrus"
)

func pulseSink(ctx *Context) {
	stash := map[pulse.ID]*BackendOptions{}

	for status := range ctx.pulseCh {
		vsID, rsID := status.Source.VsID, status.Source.RsID

		switch status.Result {
		case pulse.StatusUp:
			if opts, exists := stash[status.Source]; !exists {
				continue
			} else if _, err := ctx.UpdateBackend(vsID, rsID, opts); err != nil {
				log.Errorf("error while unstashing a backend: %s", err)
			} else {
				delete(stash, status.Source)
			}

		case pulse.StatusDown:
			if _, exists := stash[status.Source]; exists {
				continue
			} else if opts, err := ctx.UpdateBackend(vsID, rsID, &BackendOptions{Weight: 0}); err != nil {
				log.Errorf("error while stashing a backend: %s", err)
			} else {
				stash[status.Source] = opts
			}
		}

		log.Warnf("backend '%s/%s' status changed: '%s'", vsID, rsID, status.Result)
	}
}