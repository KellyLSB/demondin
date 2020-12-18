import '../utils/globalExtensions';

// React
import React from 'react';
import ReactDOM from 'react-dom';
// Apollo Client
import { ApolloProvider, client } from '../utils/apolloClient';

import Items              from './items';

import 'semantic-ui-css/semantic.min.css';

const App = () => (
  <ApolloProvider client={client}>
    <h1>Hello, world!</h1>
    <br />
    <Items />
  </ApolloProvider>
);

ReactDOM.render(<App />, document.getElementById('root'));
