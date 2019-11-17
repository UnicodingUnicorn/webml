import { GET, POST, HEAD } from './requests';
import { BASE_URL } from '../settings';

export default {
  get_batches(id) {
    return GET(`${BASE_URL}/model/${id}/batches`).then(({ body }) => body);
  },
  get_random_batch(id) {
    return GET(`${BASE_URL}/model/${id}/batch`).then(({ body }) => body);
  },
  get_data_batch(id, batchid) {
    return GET(`${BASE_URL}/model/${id}/batch/${batchid}/data`).then(({ body, headers }) => {
      return {
        id,
        parser: headers.get('x-amz-meta-parser') || '',
        shape: headers.get('x-amz-meta-shape') ? JSON.parse(headers.get('x-amz-meta-shape')) : {},
        data: body,
      };
    });
  },
  head_data_batch(id, batchid) {
    return HEAD(`${BASE_URL}/model/${id}/batch/${batchid}/data`).then(headers => {
      return {
        id,
        parser: headers.get('x-amz-meta-parser') || '',
        shape: headers.get('x-amz-meta-shape') ? JSON.parse(headers.get('x-amz-meta-shape')) : {},
      };
    });
  },
  get_label_batch(id, batchid) {
    return GET(`${BASE_URL}/model/${id}/batch/${batchid}/labels`).then(({ body, headers }) => {
      return {
        id,
        parser: headers.get('x-amz-meta-parser') || '',
        shape: headers.get('x-amz-meta-shape') ? JSON.parse(headers.get('x-amz-meta-shape')) : {},
        labels: body,
      };
    });
  },
  head_label_batch(id, batchid) {
    return HEAD(`${BASE_URL}/model/${id}/batch/${batchid}/labels`).then(headers => {
      return {
        id,
        parser: headers.get('x-amz-meta-parser') || '',
        shape: headers.get('x-amz-meta-shape') ? JSON.parse(headers.get('x-amz-meta-shape')) : {},
      };
    });
  },
  batch_data(id, dataid, options) {
    return POST(`${BASE_URL}/model/${id}/daya/${dataid}/batch`, {
      'Content-Type': 'multipart/form-data',
    }, {
      data_parser: options.data_parser,
      label_parser: options.label_parser,
      batch_size: options.batch_size,
    });
  },
};
