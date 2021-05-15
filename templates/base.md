**Bases:** `Health` - `Name` - `Affiliation`
```diff
{% for key, value in data.items() %}
{{value["diff"]}}{{value["health"]}} - {{ key }}{{value["affiliation"]}}{% endfor %}
```