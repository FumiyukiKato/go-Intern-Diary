import React from "react";
import {BrowserRouter, Route, Switch} from "react-router-dom";

import {ApolloProvider} from "react-apollo";
import ApolloClient from "apollo-client";
import {ApolloLink} from "apollo-link";
import {HttpLink} from "apollo-link-http";
import {onError} from "apollo-link-error";
import {InMemoryCache} from "apollo-cache-inmemory";

import {GlobalHeader} from "./global_header";
import {Index} from "./index";
import {Diary} from "./diary";

const client = new ApolloClient({
  link: ApolloLink.from([
    onError(({ graphQLErrors, networkError }) => {
      if (graphQLErrors)
        graphQLErrors.map(({ message, locations, path }) =>
          console.log(
            `[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`,
          ),
        );
      if (networkError) console.log(`[Network error]: ${networkError}`);
    }),
    new HttpLink({
      uri: 'http://localhost:8000/query',
      credentials: 'same-origin',
    })
  ]),
  cache: new InMemoryCache(),
});

export const App: React.StatelessComponent = () => (
  <ApolloProvider client={client}>
    <BrowserRouter basename="/spa">
      <>
        <GlobalHeader />
        <main>
          <Switch>
            <Route exact path="/" component={Index} />
            <Route exact path="/diaries/:diaryId" component={Diary} />
          </Switch>
        </main>
      </>
    </BrowserRouter>
  </ApolloProvider>
);
