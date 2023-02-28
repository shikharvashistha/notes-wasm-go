<!-- <h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation</p> -->

<script>
    var note = "";

    function isEmpty() {
        if (note == "") {
            return true;
        }
        return false;
    }

    function saveNote() {
        // TODO
        alert("TODO");
    }
</script>

<div>
    <input type="button"
        on:click={saveNote}
        value="Cloud Save">
        
    <textarea bind:value={note} placeholder="Enter your notes"></textarea>    
    <button type="submit">Submit</button>
    
    <p>: {note}</p>
</div>
