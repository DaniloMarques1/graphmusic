<script>
  import { onMount } from 'svelte';
  // because it is being exported, it will receive as props
	// export let name = "Danilo";
  const url = 'http://localhost:8080/graphql';
  let playlist = [];
  let author;
  let musicName;
  let isEdit = false;
  let oldName = '';
  let btnLabel = 'Add music';

  async function fetchMusic() {
    try {
      const response = await fetch(url, {
        method: 'post',
        body: `{findAll {id, author, name}}`,
      });
      const body = await response.json();
      const { findAll } = body;
      if (findAll) {
        playlist = findAll;
      }
    } catch(err) {
      console.log(err);
    }
  }

  async function addMusic() {
    let mutation = isEdit ?
              `mutation { updateByName(name: "${oldName}" music: { name: "${musicName}", author: "${author}" }) {id, author, name}}` :
              `mutation { addMusic(name: "${musicName}", author: "${author}") { id, author, name } }`
    try {
      const response = await fetch(url, {
        method: 'post',
        body: mutation 
      });
      const body = await response.json();
      let obj;
      if (isEdit) {
        obj = body.updateByName;
        playlist = playlist.map(m => { 
          if (m.id === obj.id) {
            m.name = obj.name;
            m.author = obj.author;
          }
          return m;
        });
        isEdit = false;
      } else {
        obj = body.addMusic;
        playlist = [...playlist, obj];
      }
      author = '';
      musicName = '';
    } catch(err) {
      console.log(err);
    }
  }

  async function removeMusic(music) {
    try {
      await fetch(url, {
        method: 'post',
        body: `mutation { removeById(id: "${music.id}") { id, author, name } }`
      });
      playlist = playlist.filter(m => m.id !== music.id);
    } catch(err) {
      console.log(err);
    }
  }

  async function updateMusic(music) {
    oldName = music.name;
    author = music.author;
    musicName = music.name;
    isEdit = true;
  }

  onMount(fetchMusic);
</script>

<main>
  <div>
    <input bind:value={author} placeholder="Author" />
  </div>
  <div>
    <input bind:value={musicName} placeholder="Music name" />
  </div>
  <button on:click={addMusic}>
    {#if isEdit}
      Edit music
    {:else}
      Add music
    {/if}
  </button>

  <ul>
    {#each playlist as music}
      <li>{music.name}<button on:click={() => removeMusic(music)}>X</button><button on:click={() => updateMusic(music)}>E</button></li>
    {/each}
  </ul>
</main>

<style>
</style>
