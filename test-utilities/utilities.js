import fetch from 'node-fetch';

export async function request(URL, accessToken, method, body) {
    const options = {
        method,
        headers: {}
    }

    if (body) {
        options.body = JSON.stringify(body);
    }

    if (accessToken) {
        options.headers["X-Florence-Token"] = accessToken;
    }

    const response = await fetch(URL, options);
    if (!response.ok) {
        throw response;
    }

    return await response.json();
}