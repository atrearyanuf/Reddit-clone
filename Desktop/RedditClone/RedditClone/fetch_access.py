import requests
import requests.auth

auth = requests.auth.HTTPBasicAuth("XshjUZznXX2tsmxZg4gVkA", "9X4p_2e8SAeRxPITxrTiGygaQgWoFA")
data = {
    'grant_type': 'password',
    'username': 'Top-Stable1774',
    'password': ''
}
headers = {'User-Agent': 'MyApp/1.0'}

response = requests.post(
    'https://www.reddit.com/api/v1/access_token',
    auth=auth,
    data=data,
    headers=headers
)
print(response.text)