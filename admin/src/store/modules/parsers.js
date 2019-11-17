import * as types from '../mutation-types';

import parserApi from '../../resource/parser';

export default {
  state: {
    parsers: [],
  },
  getters: {
    parsers: state => state.parsers,
  },
  actions: {
    init_parsers({ commit }) {
      return parserApi.get_parsers().then(parserids => {
        return Promise.all(parserids.map(parserid => {
          return parserApi.get_parser(parserid);
        })).then(parsers => {
          return parsers;
        });
      }).then(parsers => {
        commit(types.INIT_PARSERS, parsers);
      });
    },
  },
  mutations: {
    [types.INIT_PARSERS](state, parsers) {
      state.parsers = parsers;
    },
  },
}
