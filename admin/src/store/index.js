import Vue from 'vue'
import Vuex from 'vuex'

import models from './modules/models';
import parsers from './modules/parsers';

Vue.use(Vuex)

export default new Vuex.Store({
  modules: {
    models,
    parsers,
  },
  actions: {
    init({ dispatch }) {
      return dispatch('init_models').then(() => {
        return dispatch('init_parsers');
      });
    },
  },
});
