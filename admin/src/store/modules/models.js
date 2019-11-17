import * as types from '../mutation-types';

import batchApi from '../../resource/batch';
import dataApi from '../../resource/data';
import modelsApi from '../../resource/models';

export default {
  state: {
    models: [],
  },
  getters: {
    models: state => state.models,
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
  },
  mutations: {
    [types.INIT_MODELS](state, models) {
      state.models = models;
    },
  },
};
