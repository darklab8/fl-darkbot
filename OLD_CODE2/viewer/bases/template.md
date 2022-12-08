**Bases:** `Health` - `Name` - `Affiliation`
```json
{% for value in data %}
["{{value["health"]}}"]{"{{value["name"]}}"}["{{value["affiliation"]}}"]{% endfor %}
```