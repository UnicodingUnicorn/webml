import { GET, HEAD, PUT } from './requests';
import { BASE_URL } from '../settings';

export default {
  get_models () {
    return GET(`${BASE_URL}/models`).then(({ body }) => body);
  },
  get_model (id) {
    return GET(`${BASE_URL}/model/${id}`).then(({ body, headers }) => {
      return {
        id,
        name: headers.get('x-amz-meta-name') || '',
        model: body,
      };
    });
  },
  head_model (id) {
    return HEAD(`${BASE_URL}/model/${id}`).then(headers => {
      return {
        id,
        name: headers.get('x-amz-meta-name') || '',
      };
    });
  },
  put_model (id, name, data) {
    return PUT(`${BASE_URL}/model/${id}`, {
      'x-amz-meta-name': name || '',
    }, data);
  },
}
