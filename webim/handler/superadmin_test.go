package handler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/forfd/custm-chat/webim/models"
)

func TestCalculateExp(t *testing.T) {
	ent := &models.Enterprise{
		Plan:        models.EditionEnterprise,
		TrialStatus: models.TrialIn,
		AgentNum:    1,
	}

	actual := calculateExp(30, 2, ent)
	expect := time.Now().UTC().Add(time.Duration(30) * time.Hour * 24)
	assert.Equal(t, expect, actual)

	now := time.Now().UTC()
	ent = &models.Enterprise{
		Plan:           models.EditionEnterprise,
		TrialStatus:    models.TrialNone,
		AgentNum:       1,
		ExpirationTime: now.Add(time.Hour * 24),
	}

	actual = calculateExp(30, 2, ent)
	expect = now.Add(time.Duration(30)*time.Hour*24 + 12*time.Hour)
	assert.Equal(t, expect, actual)

	now = time.Now().UTC()
	ent = &models.Enterprise{
		Plan:           models.EditionEnterprise,
		TrialStatus:    models.TrialNone,
		AgentNum:       2,
		ExpirationTime: now.Add(time.Hour * 24),
	}

	actual = calculateExp(30, 2, ent)
	expect = now.Add(time.Duration(30)*time.Hour*24 + 12*time.Hour)
	assert.Equal(t, expect, actual)
}

func TestCalculateExp1(t *testing.T) {
	now := time.Now().UTC()
	ent := &models.Enterprise{
		Plan:           models.EditionEnterprise,
		TrialStatus:    models.TrialNone,
		AgentNum:       3,
		ExpirationTime: now.Add(time.Hour * 24),
	}

	actual := calculateExp(30, 3, ent)
	expect := now.Add(time.Hour*24*30 + 8*time.Hour)
	assert.Equal(t, expect, actual)
}
