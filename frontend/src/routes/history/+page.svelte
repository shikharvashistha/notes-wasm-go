<script >
    export let data
    import { P } from 'flowbite-svelte'
    import "../../app.css"
    const { history } = data 
</script>

{#each history as item}
    <P>{item.fileName} - (@{item.user},{item.date})</P>
{/each}
