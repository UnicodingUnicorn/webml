const GET = (url) => {
  return fetch(url, {
    method: 'GET',
  }).then((res) => {
    if (res.ok) {
      return ((res.headers.get('Content-Type') == 'application/json' ? res.json() : res.text()))
      .then(body => {
        return {
          body,
          headers: res.headers,
        };
      })
    } else {
      throw {
        api_error: true,
        status: res.status,
      };
    }
  });
};

const POST = (url, headers, body) => {
  return fetch(url, {
    method: 'POST',
    headers,
    body,
  }).then((res) => {
    if (res.ok) {
      if (res.headers.get('Content-Type') == 'application/json') {
        return res.json();
      } else {
        return res.text();
      }
    } else {
      throw {
        api_error: true,
        status: res.status,
      };
    }
  });
};

const PUT = (url, headers, body) => {
  return fetch(url, {
    method: 'PUT',
    headers,
    body,
  }).then((res) => {
    if (res.ok) {
      if (res.headers.get('Content-Type') == 'application/json') {
        return res.json();
      } else {
        return res.text();
      }
    } else {
      throw {
        api_error: true,
        status: res.status,
      };
    }
  });
};

const HEAD = (url, headers, body) => {
  return fetch(url, {
    method: 'HEAD',
    headers,
    body,
  }).then((res) => {
    if (res.ok) {
      return res.headers;
    } else {
      throw {
        api_error: true,
        status: res.status,
      };
    }
  });
};

export {
  GET,
  POST,
  PUT,
  HEAD,
};
