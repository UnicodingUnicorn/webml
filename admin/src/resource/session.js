import { GET, POST } from './requests';
import { BASE_URL } from '../settings';

export default {
  get_lost(id) {
    return GET(`${BASE_URL}/session/${id}/loss`).then({ body, headers } => body);
  },
  update_loss(id, loss) {
    return POST(`${BASE_URL}/session/${id}/loss`, {
      'Content-Type': 'application/json',
    }, {
      loss
    });
  },
  update_weights(id, { shape, weights }) {
    return POST(`${BASE_URL}/session/${id}/weights`, {
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
