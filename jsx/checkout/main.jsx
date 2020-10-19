import '../utils/globalExtensions';

// React
import React from 'react';
import ReactDOM from 'react-dom';
// Apollo Client
import { ApolloProvider } from '../utils/apolloClient';

// Style and Elements
import { Divider, Container, Grid, Segment } from 'semantic-ui-react';
import 'semantic-ui-css/semantic.min.css'

import Items from './items';
import Cart from './cart';

// Stripe.JS React Elements
import { Elements, StripeProvider } from 'react-stripe-elements';

ReactDOM.render(
	<ApolloProvider>
		<Container text>
			<Divider hidden />
			<Grid columns={2} divided>
				<Grid.Row stretched>
					<Grid.Column width={9}>
						<Items />
					</Grid.Column>
					<Grid.Column width={7}>
						<StripeProvider apiKey={STRIPE_PUBLISH_KEY}>
							<Elements>
								<Cart />
							</Elements>
						</StripeProvider>
					</Grid.Column>
				</Grid.Row>
				<Grid.Row stretched>
					<Grid.Column width={16} textAlign='right'>
						&copy; HeXXeD
					</Grid.Column>
				</Grid.Row>
			</Grid>
		</Container>
	</ApolloProvider>
, document.getElementById('root'));
