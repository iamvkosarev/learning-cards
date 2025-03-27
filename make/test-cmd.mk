API_URL ?= http://localhost:8080/api/learning-cards
TOKEN ?= "YOUR_TOKEN"

group_name ?= Test group
group_id ?= 1
front_text ?= こんいちは
back_text ?= hello

# Create a new group
local-test-add-group:
	@echo "Creating group: group_name=\"$(group_name)\"\n"
	curl -X POST $(API_URL)/v1/group \
	  -H "Content-Type: application/json" \
	  -H "Authorization: Bearer $(TOKEN)" \
	  -d '{"group_name":"$(group_name)"}'

# Add a card to group
local-test-add-card:
	@echo "Adding card: group_id=\"$(group_id)\", front_text =\"$(front_text)\", back_text =\"$(back_text)\""
	curl -X POST $(API_URL)/v1/card \
	  -H "Content-Type: application/json" \
	  -H "Authorization: Bearer $(TOKEN)" \
	  -d '{"group_id":$(group_id), "front_text":"$(front_text)", "back_text":"$(back_text)"}'

# Get all cards from a group
local-test-get-group-cards:
	@echo "Getting cards: group_id=\"$(group_id)\"\n"
	curl -X GET $(API_URL)/v1/group/$(group_id)/cards \
	  -H "Authorization: Bearer $(TOKEN)"