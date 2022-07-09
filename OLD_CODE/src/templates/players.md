**{{ title }}** `Name` - `System` - `Online` {% if alert == true %} alert level above the threshold, alert @here {% endif %}
```json
{% for key, value in data.items() %}
["{{value["name"]}}"]{"{{ value["system"] }}"}["{{value["time"]}}"]{% endfor %}
```

