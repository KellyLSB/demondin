import React from 'react'

export default class Price extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            data: 'data' in props ? props.data : {}
        };

        this.onChangePrice = this.onChangePrice.bind(this)
        this.onChangeValidAfter = this.onChangeValidAfter.bind(this)
        this.onChangeValidBeofre = this.onChangeValidBefore.bind(this)
    }

    onChangePrice(event) {
        var value = event.target.value
        this.setState((state) => {
            state.data.Price = value
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

    }

    render() {
        return (
            <form>
                <input placeholder="$0.00"
                    onChange={this.onChangePrice}
                    value={this.state.data.Price} />
                <input placeholder="Valid From"
                    onChange={this.onChangeValidAfter}
                    value={this.state.data.ValidAfter} />
                <input placeholder="Valid Until"
                    onChange={this.onChangeValidBefore}
                    value={this.state.data.ValidBefore} />
            </form>
        )
    }
}
