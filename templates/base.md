**Bases:**
```json
{% for key, value in data.items() %}
["{{value["health"]}}"]{"{{ key }}"}["{{value["affiliation"]}}"]{% endfor %}
```