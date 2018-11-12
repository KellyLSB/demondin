import React from 'react'
import NumberFormat from 'react-number-format';


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
        this.onChangeValidBeofre = this.onChangeValidBefore.bind(this)
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
        var value = event.target.value
        this.setState((state) => {
            state.data.AfterDate = value
            return state
        })
    }

    onChangeValidBefore(event) {
        var value = event.target.value
        this.setState((state) => {
            state.data.ValidBefore = value
            return state
        })
    }

    onSubmit(event) {
        event.preventDefault()
        // Sanitize data before submitting to the API
        
        alert("submit prices")
        console.log(this.state.data)
        
        fetch(`/shop/keeper/badges/${this.props.item}/prices/${this.state.data.ID}.json`, {
            method: "ID" in this.state.data ? "PUT" : "POST",
            body: JSON.stringify(this.state.data)
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
                <input placeholder="Valid From"
                    value={this.state.data.ValidAfter}
                    onChange={this.onChangeValidAfter} />
                <input placeholder="Valid Until"
                    value={this.state.data.ValidBefore}
                    onChange={this.onChangeValidBefore} />
                <input type="submit" value="Save" />
            </form>
        )
    }
}
