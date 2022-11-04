/**
 * Handler that will be called during the execution of a PostLogin flow.
 *
 * @param {Event} event - Details about the user and the context in which they are logging in.
 * @param {PostLoginAPI} api - Interface whose methods can be used to change the behavior of the login.
 */

exports.onExecutePostLogin = async (event, api) => {
  const axios = require("axios");

  if (!("uuid" in event.user.app_metadata)) {
    const url = "https://downloop.us.auth0.com/oauth/token";
    const body = {
      client_id: event.secrets.AUTH0_CLIENT_ID,
      client_secret: event.secrets.AUTH0_CLIENT_SECRET,
      audience: "https://api.downloop.io",
      grant_type: "client_credentials",
    };

    const token = await axios
      .post(url, body)
      .then(function (response) {
        return response.data.access_token;
      })
      .catch(function (error) {
        console.log("Failed to get access token: " + error);
        throw new Error(error);
      });

    const downloopURL = "https://downloop.tunnelto.dev/users";
    const downloopBody = {
      username: event.user.email,
    };

    const uuid = await axios
      .post(downloopURL, downloopBody, {
        headers: {
          Authorization: "Bearer " + token,
          "Content-Type": "application/json",
        },
      })
      .then(function (response) {
        return response.data.data.id;
      })
      .catch(function (error) {
        console.log("Failed to post user to downloop api: " + error);
        throw new Error(error);
      });

    api.user.setAppMetadata("uuid", uuid);
  }
};

/**
 * Handler that will be invoked when this action is resuming after an external redirect. If your
 * onExecutePostLogin function does not perform a redirect, this function can be safely ignored.
 *
 * @param {Event} event - Details about the user and the context in which they are logging in.
 * @param {PostLoginAPI} api - Interface whose methods can be used to change the behavior of the login.
 */
// exports.onContinuePostLogin = async (event, api) => {
// };
