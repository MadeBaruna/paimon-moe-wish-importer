<script>
  import { onMount } from "svelte";

  import "./global.css";
  import Main from "./main.svelte";

  onMount(() => {
    if (document.body.scrollHeight > document.body.clientHeight) {
      window.go.main.App.ResizeWindow(600, document.body.scrollHeight + 160);
    }
  });
</script>

<Main />
