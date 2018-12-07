import React from 'react'
import { Grid, Header, Segment, Button, Icon, Label } from 'semantic-ui-react'


export default class Item extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: {}
    };

    if ('data' in props) {
      this.state.data = props.data;
    }
  }
  
  toDollars(cents) {
    return (cents / 100).toLocaleString("en-US", {
      style:"currency", currency:"USD"
    })
  }
    
  toCents(dollars) {
    return Math.round(100 * parseFloat(
      `${dollars}`.replace(/[$,]/g, '')
    ))
  }
  
  currentPrice() {
    return this.state.data.Prices[0]
  }
  
  printPrice() {
    var price = this.currentPrice()
    
    if(price.Taxable) {
      return this.toDollars(price.Price) + " + Tax"
    }
    
    return this.toDollars(price.Price)
  }

  render() {
    return (
      <Segment>
        <Header size="huge">{this.state.data.Name}</Header>
        <span>{this.state.data.Description}</span>
        <Grid columns={2}>
          <Grid.Column>
            <Header sub>Price</Header>
            <Label tag>{this.printPrice()}</Label>
          </Grid.Column>
          <Grid.Column textAlign='right'>
            <Button
              primary icon labelPosition='left'
              onClick={this.props.addToCart(this.state.data.ID)}>
              <Icon name='shop' />
              Add To Cart
            </Button>
          </Grid.Column>
        </Grid>
      </Segment>
    );
  }
}
