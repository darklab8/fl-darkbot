jq -n \
--arg hetzner_token $(pass api/personal/terraform/hetzner/production) \
--arg cloudflare_token $(pass api/personal/terraform/cloudflare/dd84ai) \
'{
    "hetzner_token": $hetzner_token,
    "cloudflare_token": $cloudflare_token
}'
