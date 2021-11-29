<script>
  import { onMount } from "svelte";

  const types = {
    100: "Beginners' Wish",
    200: "Standard",
    301: "Character Event",
    302: "Weapon Event",
  };

  const total = {};

  let error = "";
  let server = "global";
  let step = 0;
  let page = 1;
  let uid = "";
  let currentBanner = "Beginners' Wish";
  let result = "";
  let copied = false;
  let saved = false;
  let copyText = "COPY";
  let textarea;

  function changeServer(val) {
    server = val;
  }

  async function start() {
    copied = false;
    saved = false;
    uid = "";
    error = "";
    currentBanner = "Beginners' Wish";
    result = "";

    step = 1;
    await window.go.main.App.Start(server);
  }

  function cancel() {
    window.go.main.App.Cancel();
    step = 0;
  }

  function back() {
    step = 0;
  }

  function copy() {
    window.go.main.App.Copy(result);

    copyText = "COPIED";
    setTimeout(() => {
      copyText = "COPY";
    }, 2000);
  }

  function focusSelectAll() {
    textarea.focus();
    textarea.select();
  }

  onMount(() => {
    window.runtime.EventsOn("error", (err) => {
      console.log(err);
      step = -1;
      switch (err.name) {
        case "LOGNOTFOUND":
          error =
            "Cannot find the log file! Make sure to open the wish history first!";
          break;
        case "LINKNOTFOUND":
          error =
            "Cannot find the wish history link! Make sure to open the wish history first!";
          break;
        case "ERRORREADLOG":
          error = "Error when reading the log file! " + err.message;
          break;
      }
    });

    window.runtime.EventsOn("uid", (val) => {
      uid = val;
    });

    window.runtime.EventsOn("page", (val) => {
      page = val;
    });

    window.runtime.EventsOn("banner", (val) => {
      currentBanner = types[val];
    });

    window.runtime.EventsOn("total", (val) => {
      if (total[val.code] === undefined) total[val.code] = 0;
      total[val.code] = val.total;
    });

    window.runtime.EventsOn("result", (val) => {
      result = val;
      step = 2;
    });

    window.runtime.EventsOn("copied", () => {
      copied = true;
    });

    window.runtime.EventsOn("saved", () => {
      saved = true;
    });
  });
</script>

<main class="p-4 flex flex-col h-full select-none" data-wails-no-drag>
  <div class="flex items-center">
    <h1 class="text-primary font-bold text-xl font-display flex-1">
      Paimon.moe Wish Importer
    </h1>
    <p class="text-gray-300">
      {#if uid !== ""}
        <span class="mr-4">UID {uid}</span>
      {/if} v1
    </p>
  </div>
  {#if step === 0}
    <h1 class="text-white font-bold text-3xl font-display">
      Choose your server
    </h1>
    <div class="flex mt-2">
      <button
        class="pill {server === 'global' ? 'active' : ''}"
        on:click={() => changeServer("global")}>Global</button
      >
      <button
        class="pill {server === 'china' ? 'active' : ''}"
        on:click={() => changeServer("china")}>China</button
      >
    </div>
    <div class="flex-1" />
    <ol class="list-decimal list-inside">
      <li class="font-body text-white">Open Genshin Impact on this PC</li>
      <li class="font-body text-white">Open the wish history (in the game)</li>
      <li class="font-body text-white">Press start import below!</li>
    </ol>
    <button
      class="ring-2 ring-primary rounded-xl px-4 py-2 mt-4 text-primary font-body font-bold hover:bg-primary hover:text-white transition duration-100"
      on:click={start}
    >
      START IMPORT
    </button>
  {:else if step >= 1}
    {#if step === 1}
      <div class="flex items-center mt-2">
        <div class="spinner mr-4" />
        <div class="flex flex-col text-white font-body flex-1">
          <p>Processing {currentBanner} Banner</p>
          <p>Page {page}</p>
        </div>
        <button
          class="ring-2 ring-gray-500 rounded-xl px-2 py-1 text-red-400 text-sm
        hover:ring-red-400 hover:bg-red-400 hover:text-white transition duration-100 font-body"
          on:click={cancel}
        >
          Cancel
        </button>
      </div>
      <div class="flex-1" />
    {/if}
    <div
      class="{Object.keys(total).length === 0
        ? ''
        : 'border'} border-gray-700 rounded-xl min-w-full mt-2"
    >
      <table class="progress w-full font-body">
        {#each Object.entries(total) as [code, total]}
          <tr>
            <td class="px-2 py-1">
              <span class="text-white mr-2 whitespace-no-wrap">
                {types[code]} Banner
              </span>
            </td>
            <td class="pr-2 py-1 text-right">
              <span class="text-white mr-2 whitespace-no-wrap">
                x {total}
              </span>
            </td>
          </tr>
        {/each}
      </table>
    </div>
    {#if step === 2}
      <div class="flex-1" />
      <div class="flex font-body">
        <div
          class="rounded-xl overflow-hidden border-secondary focus-within:border-primary border-2 border-transparent ease-in duration-100 flex-1 mr-4"
        >
          <textarea
            bind:this={textarea}
            on:click={focusSelectAll}
            class="bg-transparent bg-secondary text-xs text-gray-200 px-2 w-full h-full resize-none overflow-hidden select-all"
            >{result}</textarea
          >
        </div>
        <button
          class="border-2 border-primary rounded-xl px-4 py-2 text-primary font-body font-bold hover:bg-primary hover:text-white transition duration-100 h-full"
          on:click={copy}
        >
          {copyText}
        </button>
      </div>
      <p class="text-xs text-gray-400 ml-2 font-body">
        {#if copied}<span>Copied to clipboard.</span>{/if}
        {#if saved}<span> Saved to Downloads folder.</span>{/if}
      </p>
      <p class="text-xs text-primary ml-2 -mb-2 font-body">
        Paste it back to paimon.moe!
      </p>
    {/if}
  {:else if step === -1}
    <h1 class="text-white font-bold text-3xl font-display">Error</h1>
    <p class="font-body text-red-400 flex-1 break-words select-text">{error}</p>
    <button
      class="ring-2 ring-gray-500 rounded-xl px-4 py-2 mt-4 text-white font-body font-bold hover:bg-primary hover:ring-primary transition duration-100"
      on:click={back}
    >
      Back
    </button>
  {/if}
</main>

<style lang="postcss">
  .pill {
    @apply font-body rounded-2xl ring-2 ring-gray-500 px-4 py-1 mr-3 text-white transition duration-100;
    @apply hover:bg-primary hover:ring-primary;
  }

  .active {
    @apply bg-primary ring-primary;
  }

  .spinner {
    border: 4px rgba(255, 255, 255, 0.25) solid;
    border-top: 4px white solid;
    width: 32px;
    height: 32px;
    border-radius: 50%;
    animation: spin 2s infinite linear;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(359deg);
    }
  }

  table.progress tr:not(:last-child) {
    @apply border-b border-gray-700;
  }
</style>
