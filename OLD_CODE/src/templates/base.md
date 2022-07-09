**Bases:** `Health` - `Name` - `Affiliation`
```json
{% for key, value in data.items() %}
{"{{ value["diff"] }}"}["{{value["health"]}}"]{"{{ key }}"}["{{value["affiliation"]}}"]{% endfor %}
```