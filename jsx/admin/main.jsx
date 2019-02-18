import React              from 'react';
import ReactDOM           from 'react-dom';
import ApolloClient       from "apollo-boost";
import { ApolloProvider } from "react-apollo";
import Items              from './items';

import '../../semantic/dist/semantic.min.css';

const client = new ApolloClient({
  uri: "http://localhost:8080/"
});

const App = () => (
  <ApolloProvider client={client}>
    <h1>Hello, world!</h1>
    <br />
    <Items />
  </ApolloProvider>
);

ReactDOM.render(<App />, document.getElementById('root'));
