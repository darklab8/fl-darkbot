  Amount({{ .AmountValue |  printf "%-06d" }}), {{ capitalize .Category }}({{ .GoodName | printf "%q" }}){{ if .IsEnd }};{{end}}
