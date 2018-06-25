import fetch from 'node-fetch';

export async function request(URL, accessToken, method, body) {
    const options = {
        body: JSON.stringify(body),
        method,
        headers: {}
    }

    if (accessToken) {
        options.headers = accessToken
    }

    const response = await fetch(URL, options);
    if (!response.ok) {
        throw response;
    }

    return await response.json();
}