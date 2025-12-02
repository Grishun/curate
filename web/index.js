const {createApp, ref} = Vue

createApp({
    data: function () {
        return {
            service_name: "curate",
            rates: [],
        }
    },

    methods: {
        getRates: function () {
            axios.get('http://localhost:8080/api/v1/rates')
                .then(function (response) {
                    console.log(response)
                    this.rates = response.data
                }).catch(function (error) {
                    console.log(error)
                })
        },

    },

    mounted: function () {
        this.getRates()
    },

}).mount("#app")