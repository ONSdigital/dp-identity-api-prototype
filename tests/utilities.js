import fetch from 'node-fetch';

export default class Utilities {
    static async request(URL) {
        const options = {

        }
        await fetch(URL, options)
    }
}