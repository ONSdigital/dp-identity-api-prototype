import fs from 'fs-extra';

import {request} from '../test-utilities/utilities';

const userJSON = {"passwordHash":"pyoLEXCvzEiY218C9VTHRzPbOHqgghLsPhftp9qLvRRx8t15ERA93zgnpV/vtR87","keyring":{"privateKeySalt":"b5l/tQ6oTKFuyXMACC37Mg\u003d\u003d","privateKey":"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\u003d","publicKey":"MIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAtAjEv1Us3dnHW0YJfwwlc/M53ZHyy7Bq1HZRE+KsAScH9Z01ImTeWUhdRDwanhn3+XJWKRjZ7TdqNCetkH/qfyGOjzTc5ohvrdiUgfVhAOhD0wJm25ARd/vq68lSCW7tqKuD0HUfzydBPYtheQLX9+nS/UJMAAFW1dS1Nk3R5ytuQMaaaXkk7mt/BBaGOKaG2o6yTWEiS7tPfbrxfpNl++c9n3JwnSX0yUbAFwOG98en8mmluN3uBaDvapMfcmi2B8M2A5u4BVQ3t6vF/RGQHdtolvqEHANblLJg5b0w/MAN7BrIa08e2mydP7R9am2p+A0pUjYuT1fQZ5OrpV75Mcjaiylou6EvcEJsjP/q3Ke4SAuSRBKnPV7IBkGaDvNc/p+/cV2mBPX32dxt7HikEHmf6/jJBbnij2AFnjgoi/eZhaOUTqKc+9aZorx6FXig6585WC/jwjrVb8mo4VZDieZQ4Upv43L+xQ+j5nCGp1gAA8VkbFQPtDcJWbN4iBxPAgMBAAE\u003d","keyring":{"asdsadasdasd-fe8fcea2051b69791a7063f8f8b3893023d0d76de605b9e833dad47a8a7e1efc":"VaMMgE95dBE7NY6h5I6F7wn0BrVw2dNP5zdn7AF5ZnYQNizywU1PM/8qu6yVYJpr91fJZUVBN568zb6VYAym4NYBnLuIZu6JImjAr9enBgAV9CuPBl7EX8dYArGCCTsMMLE2RYh5ylQ7v2FdqFa29Efd5oUjZtdtpnWnitxZ7AgkA5EykSQCwHGfXLVpQKyGRLH3YNkgAV0UCtulSTG+3AXwbG44/sw6NlvijA9FBGX8CgGJElK4Ebi0pmFvyKdxVDUdm7SrvMd/KSsaxuj/Br3hARpp5u4fHgtdWF4s+EP1aYZpJN9AOwWLDGy9GaMoyKD7QOtdy80gN7f69I/7fkpHPLxQXICP0Toekz8+5bcbHATDx1sDkMfun55RAJBEV8dYPjxHQUrjBSwJDN08d0uP6vPjCjhYJqduuIMrvvKMPOiP0gDMF+EveE0zneRKLh4bam4ynCSnjAEcGx2X4hrRL3eu8jxGzQ2k6iTkEgei1LjDkrpVimPeOp7yPcwI","csdb-import":"fRz5FrvQABkAR1ThCYkg+5jMbIbKIbTe3+t0NkUjJQ2qTsVO1/1OLgUauSXa6zXc6t4bExgLyDMR0WeTefxatl76v35Jg/Da4EK03S7rY7KWQa3I2WJvmI3X7XV8cEyLN3CxbHCkd4iVSSYGfjfF5GqCUjU1nQqhYXTiWVoO81S+e459o6VCpCklrwp2AJvwXDjfuTwDs55tlJjqaHoZkxJUhxVdI/3nEJq98qY7lgm4Jn1HiydapYu4KsDSphDaCJtgeYwlhqkz2ElHI0eOGpuskuUBUMiEdgYtCe1b5WAX/DSzhclrNUjjb68vKua5KjfaM9vGuxO2i8TSFXJZWN6ehqo1b9H/PQg7hv7ABA5ChiQ4WI3IuUwEqjL9yJp3RCry5fK19pGBfQ1v25ZFjx2U4iVlRWsTnEFPiis2XibIm6g83cvYLA1V0s8ExSZ/7X/acHzjACa/U1LTFDyVaQ/KkhXu1rS9Mr/Pbj0A6nc8jqkNOUw4MdyczVBRYLSD"}},"name":"Identity API test user","email":"test-identity-api@email.com","inactive":false,"temporaryPassword":false,"lastAdmin":"test-identity-api@email.com"};

const accessToken = '12345-abcde-67890-fghij';

const config = {
    identityAPIUrl: process.env.IDENTITY_API_URL || "http://localhost:23700",
    zebedeeUrl: process.env.ZEBEDEE_URL || "http://localhost:8082",
    zebedeeDir: `${process.env.zebedee_root}/zebedee`,
};

const writeUserToDisk = async json => {
    try {
        console.log(`Writing test user data to '${config.zebedeeDir}/users/testidentityapiemail.com.json'`);
        await fs.writeJSON(`${config.zebedeeDir}/users/testidentityapiemail.com.json`, json);
    } catch (error) {
        console.error("Error writing user to disk", error);
        process.exit(1);
    }
};

beforeAll(async () => {
    await writeUserToDisk(userJSON);
});

describe("Unmigrated user login", () => {
    
    it("Identity API returns 303 error on first attempt to login with valid email", async () => {
        const body = {
            email: "test-identity-api@email.com",
            password: "one two three four"
        };
        try {
            await request(`${config.identityAPIUrl}/login`, null, "POST", body);
            fail("Expected 303 response from Identity API but received 200");
        } catch (error) {
            expect(error.status).toBe(303);
        }
    });
    
    it("Identity API returns 401 error on attempt to login with invalid email", async () => {
        const body = {
            email: "incorrect@email.com",
            password: "one two three four"
        };
        try {
            await request(`${config.identityAPIUrl}/login`, null, "POST", body);
            fail("Expected 401 response from Identity API but received 200");
        } catch (error) {
            expect(error.status).toBe(401);
        }
    });
    
    it("Identity API returns 401 error on attempt to login with invalid password", async () => {
        const body = {
            email: "test-identity-api@email.com",
            password: "wrong password"
        };
        try {
            await request(`${config.identityAPIUrl}/login`, null, "POST", body);
            fail("Expected 401 response from Identity API but received 200");
        } catch (error) {
            expect(error.status).toBe(401);
        }
    });

    it("Zebedee returns 200 on first attempt to login after migration", async () => {
        const body = {
            email: "test-identity-api@email.com",
            password: "one two three four"
        };
        try {
            const response = await request(`${config.zebedeeUrl}/login`, null, "POST", body);
            expect(response).toBe("User has been migrated");
        } catch (error) {
            fail("Expected 200 response from Zebedee");
        }

    });

    it("Zebedee returns error when login in as migrated user", async () => {
        const body = {
            email: "test-identity-api@email.com",
            password: "one two three four"
        };
        try {
            const response = await request(`${config.zebedeeUrl}/login`, null, "POST", body);
            fail("Expected 401 response from Zebedee but received 200");
        } catch (error) {
            expect(error.status).toBe(401);
        }
    });
    
    it.skip("Identity API returns access token on attempt after user's successful login to Zebedee", () => {

    });

    it("Requests to Zebedee with access token return 200 and valid response body", async () => {
        try {
            const collections = await request(`http://localhost:8081/zebedee/collections`, accessToken, "GET", null);
            expect(collections.length).toBeGreaterThan(0);
        } catch (error) {
            fail(error);
        }
    });
});