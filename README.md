# Cloudflare AntiBot Browser Challenge
This repository contains research about the new Cloudflares anti-bot javascript challenge.
This repository is for educational purposes only!

## Javascript Challenge
Bellow is a list of steps required to solve the js challenge:
1. Send a `GET` request to `cdn-cgi/challenge-platform/orchestrate/jsch/v1`, which replies with a obfuscated javascript file to generate the challenge id and make the next request.
2. Send a `POST` request to `cdn-cgi/challenge-platform/generate/ov1/{solved-challenge-id}/{cloudflare-ray-id}/cf_chl_1 cookie-here`. This request should contain the following cookies:
    - `__cfuid`: Cloudflare request id
    - `cf_chl_1`: Cloudflare challenge 1 id

3. Send a `POST` request like the second request, but with the `cf_chl_rc_ni` cookie.
4. Final request is also a `POST` request to `?__cf_chl_jschl_tk__={generated-token}`, unless the follow up is a captcha challenge.

### Generate Challenge ID
In step 1 a javascript file is aquired, this file can be used to generate the challenge id.
...describe challenge id generation...

## Captcha Challenge
I'm yet to find out how to solve this challenge.