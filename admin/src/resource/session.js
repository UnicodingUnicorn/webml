import { GET, POST } from './requests';
import { BASE_URL } from '../settings';

export default {
  get_sessions(model_id) {
    return GET(`${BASE_URL}/model/${model_id}/sessions`).then({ body, header } => body);
  },

  get_session(model_id, session_id) {
    return GET(`${BASE_URL}/model/${model_id}/session/${session_id}`).then({ body, header } => body);
  },

  update_loss(model_id, session_id, loss) {
    return POST(`${BASE_URL}/model/${model_id}/session/${session_id}/loss`, {
      'Content-Type': 'application/json',
    }, {
      loss
    });
  },
  update_weights(model_id, session_id, { shape, weights }) {
    return POST(`${BASE_URL}/model/${model_id}/session/${session_id}/weights`, {
      'Content-Type': 'application/json',
    }, {
      shape,
      data: weights,
    });
  },
  new_session(id, { shape, loss, alpha }) {
    return POST(`${BASE_URL}/session/${id}`, {
      'Content-Type': 'application/json',
    }, {
      shape,
      loss,
      alpha,
    });
  },
};
