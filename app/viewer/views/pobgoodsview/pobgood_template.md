Base({{ .BaseName | printf "%q" }}), Amount({{ .AmountValue |  printf "%-06d" }}), {{ capitalize .Category }}({{ .GoodName | printf "%q" }})
