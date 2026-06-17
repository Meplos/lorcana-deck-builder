package domain

type (
	CardAbility struct {
		Name        string `bson:"name"`
		Description string `bson:"description"`
	}
	CollectionCard struct {
		ID              string        `bson:"_id"`
		Name            string        `bson:"name"`
		Set             string        `bson:"set"`
		Version         int           `bson:"ver"`
		Title           string        `bson:"title"`
		Cost            int           `bson:"cost"`
		Type            string        `bson:"type"`
		Number          string        `bson:"number"`
		Color           []InkColor    `bson:"color"`
		Illustrator     string        `bson:"illustrator"`
		Lore            int           `bson:"lore"`
		Strength        int           `bson:"strength"`
		Willpower       int           `bson:"willpower"`
		Movement        int           `bson:"movement"`
		Ink             int           `bson:"ink"`
		Characteristics string        `bson:"characteristics"`
		Abilities       []CardAbility `bson:"abilities"`
		Variants        string        `bson:"variants"`
		Rarity          string        `bson:"rarity"`
		Language        string        `bson:"language"`
		Path            string        `bson:"path"`
		Franchise       string        `bson:"franchise"`
		Ordinal         int           `bson:"ordinal"`
		Formats         string        `bson:"formats"`
		Quantity        int           `bson:"quantity"`
	}
	Collection struct {
		ID    string                    `bson:"_id,omitempty"`
		Name  string                    `bson:"name"`
		Cards map[string]CollectionCard `bson:"cards"`
	}
)
