import { GET, PUT, HEAD } from './requests';
import { BASE_URL } from '../settings';

export default {
  get_model_data (id) {
    return GET(`${BASE_URL}/model/${id}/data`).then(({ body }) => body);
  },
  get_model_data_by_id (id, dataid) {
    return GET(`${BASE_URL}/model/${id}/data/${dataid}`).then(({ body, headers }) => {
      return {
        id,
        parser: headers.get('x-amz-meta-parser') || '',
        shape: headers.get('x-amz-meta-shape') ? JSON.parse(headers.get('x-amz-meta-shape')) : {},
        data: body,
      };
    });
  },
  head_model_data_by_id (id, dataid) {
    return HEAD(`${BASE_URL}/model/${id}/data/${dataid}`).then(headers => {
      return {
        id,
        parser: headers.get('x-amz-meta-parser') || '',
        shape: headers.get('x-amz-meta-shape') ? JSON.parse(headers.get('x-amz-meta-shape')) : {},
      };
    });
  },
  put_model_data (modelid, dataid, data, shape, parserid) {
    return PUT(`${BASE_URL}/model/${modelid}/data/${dataid}`, {
      'Content-Type': 'multipart/form-data',
      'x-amz-meta-shape': shape ? JSON.stringify(shape) : '',
      'x-amz-meta-parser': parserid || '',
    }, data);
  },
  get_model_labels (id) {
    return GET(`${BASE_URL}/model/${id}/labels`).then(({ body }) => body);
  },
  get_model_label_by_id (id, labelid) {
    return GET(`${BASE_URL}/model/${id}/label/${labelid}`).then(({ body, headers }) => {
      return {
        id,
        parser: headers.get('x-amz-meta-parser') || '',
        shape: headers.get('x-amz-meta-shape') ? JSON.parse(headers.get('x-amz-meta-shape')) : {},
        labels: body,
      };
    });
  },
  head_model_label_by_id (id, labelid) {
    return HEAD(`${BASE_URL}/model/${id}/label/${labelid}`).then(headers => {
      return {
        id,
        parser: headers.get('x-amz-meta-parser') || '',
        shape: headers.get('x-amz-meta-shape') ? JSON.parse(headers.get('x-amz-meta-shape')) : {},
      };
    });
  },
  put_model_label (modelid, labelid, labels, shape, parserid) {
    return PUT(`${BASE_URL}/model/${modelid}/labels/${labelid}`, {
      'Content-Type': 'multipart/form-data',
      'x-amz-meta-shape': shape ? JSON.stringify(shape) : '',
      'x-amz-meta-parser': parserid || '',
    }, labels);
  },
};
