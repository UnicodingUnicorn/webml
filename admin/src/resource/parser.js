import { GET, HEAD, PUT } from './requests';
import { BASE_URL } from '../settings';

export default {
  get_parsers() {
    return GET(`${BASE_URL}/parsers`).then(({ body }) => body);
  },
  get_parser(id) {
    return GET(`${BASE_URL}/parser/${id}`).then(({ body, headers }) => {
      return {
        id,
        name: headers.get('x-amz-meta-name') || '',
        parser: body,
      };
    });
  },
  head_parser(id) {
    return HEAD(`${BASE_URL}/parser/${id}`).then(headers => {
      return {
        id,
        name: headers.get('x-amz-meta-name') || '',
      };
    });
  },
  put_parser(data, name) {
    return PUT(`${BASE_URL}/parser`, {
      'Content-Type': 'multipart/form-data',
      'x-amz-meta-name': name || '',
    }, data);
  },
};
