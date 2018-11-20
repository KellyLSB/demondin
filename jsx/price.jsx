import React from 'react'
import NumberFormat from 'react-number-format';
import DatePicker from 'react-datepicker';

import 'react-datepicker/dist/react-datepicker.css';

export default class Price extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      data: 'data' in props ? props.data : {}
    };

    this.onSubmit = this.onSubmit.bind(this)
    this.onChangePrice = this.onChangePrice.bind(this)
    this.onValueChangePrice = this.onValueChangePrice.bind(this)
    this.onChangeValidAfter = this.onChangeValidAfter.bind(this)
    this.onChangeValidBefore = this.onChangeValidBefore.bind(this)
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

  onChangePrice(event) {
    var value = this.toCents(event.target.value)
    this.setState((state) => {
      state.data.Price = value
      return state
    })
  }
    
  onValueChangePrice(values) {
    const {formattedValue, value} = values
    this.setState((state) => {
      state.data.Price = this.toCents(value)
      return state
    })
  }

  onChangeValidAfter(event) {
    var value = 'target' in event ? event.target.value: event
    this.setState((state) => {
      state.data.ValidAfter = value
      return state
    })
  }

  onChangeValidBefore(event) {
    var value = 'target' in event ? event.target.value: event
    this.setState((state) => {
      state.data.ValidBefore = value
      return state
    })
  }

  onSubmit(event) {
    event.preventDefault()
    
    if ('onSubmit' in this.props) {
      this.props.onSubmit(this.state.data)
      return
    }
        
    return this.onSubmitXHR(this.state.data)
  }
    
  onSubmitXHR(xhrData) {
    var page = 'ID' in xhrData ? `/${xhrData.ID}` : ''
    fetch(`/shop/keeper/badges/${this.props.item}/prices${page}.json`, {
      method: "ID" in xhrData ? "PUT" : "POST",
      body: JSON.stringify(xhrData)
    }).then((response) => response.json()).then((data) => {
      this.setState((state) => {
        state.data = data
        return state
      })
    })
  }

  render() {
    return (
      <form onSubmit={this.onSubmit}>
        <NumberFormat thousandSeparator={true} prefix={'$'} decimalScale={2}
            value={this.toDollars(this.state.data.Price)}
            onValueChange={this.onValueChangePrice} />
            
        <DatePicker selectsStart
          selected={this.state.data.ValidAfter}
          startDate={this.state.data.ValidAfter}
          endDate={this.state.data.ValidBefore}
          onChange={this.onChangeValidAfter} />
        <DatePicker selectsEnd
          selected={this.state.data.ValidBefore}
          startDate={this.state.data.ValidAfter}
          endDate={this.state.data.ValidBefore}
          onChange={this.onChangeValidBefore} />
            
        <input type="submit" value="Save" />
      </form>
    )
  }
}
