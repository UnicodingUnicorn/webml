import * as types from '../mutation-types';

import crypto from 'crypto';

import batchApi from '../../resource/batch';
import dataApi from '../../resource/data';
import modelsApi from '../../resource/models';

export default {
  state: {
    models: [],
  },
  getters: {
    models: state => state.models,
    model: state => id => state.models.find(m => m.id === id) || null,
  },
  actions: {
    init_models({ commit }) {
      return modelsApi.get_models().then(models => {
        return Promise.all(models.map(modelid => {
          return modelsApi.get_model(modelid).then(model => {
            return dataApi.get_model_data(modelid).then(dataids => {
              return Promise.all(dataids.map(dataid => dataApi.head_model_data_by_id (modelid, dataid)));
            }).then(data => {
              model.data = data;
              return dataApi.get_model_labels(modelid);
            }).then(labelids => {
              return Promise.all(labelids.map(labelid => dataApi.head_model_label_by_id(modelid, labelid)));
            }).then(labels => {
              model.labels = labels;
              return batchApi.get_batches(modelid);
            }).then(batches => {
              model.batches = batches;
              return model;
            });
          }).then(model => {
            return model;
          });
        }));
      }).then(models => {
        commit(types.INIT_MODELS, models);
      });
    },
    add_model({ commit }, { name, model }) {
      const id = crypto.randomBytes(16).toString('hex');
      return modelsApi.put_model(id, name, model).then(() => {
        const model = {
          id,
          name,
          model,
          data: [],
          labels: [],
          batches: [],
        };
        commit(types.ADD_MODEL, model);
      }).catch(e => {
        if (e == 400) {
          throw 'Incorrect data supplied';
        } else if (e == 409) {
          throw 'ID conflict';
        } else if (e == 500) {
          throw 'Server error encountered';
        } else if (typeof e === 'string') {
          throw e;
        } else {
          throw 'A client error has occurred';
        }
      });
    },
    add_dataset({ commit, getters, rootGetters }, { model, data, labels, shape, dataparser, labelparser }) {
      return (new Promise((resolve, reject) => {
        try {
          const parsedShape = JSON.parse(shape);
          resolve(parsedShape);
        } catch(e) {
          reject('invalid shape JSON');
        }
      })).then(parsedShape => {
        if (rootGetters.parser(dataparser) == null) {
          throw 'parser not found';
        } else if (rootGetters.parser(labelparser) == null) {
          throw 'parser not found';
        } else if (getters.model(model) == null) {
          throw 'model not found';
        }
        return parsedShape;
      }).then(parsedShape => {
        const id = crypto.randomBytes(16).toString('hex');
        return dataApi.put_model_data(model, id, data, parsedShape, dataparser).then(() => {
          return dataApi.put_model_labels(model, id, labels, parsedShape, labelparser);
        }).then(() => {
          return dataApi.head_model_data_by_id(model, `data:${id}`);
        }).then(data => {
          commit(types.ADD_MODEL_DATA, { modelid: model, data });
        }).then(() => {
          return dataApi.head_model_label_by_id(model, `labels:${id}`);
        }).then(labels => {
          commit(types.ADD_MODEL_LABELS, { modelid: model, labels });
        });
      }).catch(e => {
        if (e == 400) {
          throw 'incorrect data supplied';
        } else if (e == 404) {
          throw 'model not found';
        } else if (e == 500) {
          throw 'server error';
        } else if (typeof e === 'string') {
          throw e;
        } else {
          throw 'client error';
        }
      });
    },
  },
  mutations: {
    [types.INIT_MODELS](state, models) {
      state.models = models;
    },
    [types.ADD_MODEL](state, model) {
      state.models.push(model);
    },
    [types.ADD_MODEL_DATA](state, { modelid, data }) {
      const model = state.models.find(m => m.id == modelid);
      if (model != undefined) {
        model.data.push(data);
      }
    },
    [types.ADD_MODEL_LABELS](state, { modelid, labels }) {
      const model = state.models.find(m => m.id == modelid);
      if (model != undefined) {
        model.labels.push(labels);
      }
    },
  },
};
