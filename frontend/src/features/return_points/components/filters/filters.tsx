import Styles from "./filters.module.sass";
import Button from "@components/button/button";
import IconButton from "@components/icon_button/icon-button";

import { useRef, useState } from "preact/compat";

import type { JSX } from "preact/jsx-runtime";

const FILTER_OPTIONS = [
  {
    label: "Wszystkie",
    value: "all",
  },
  {
    label: "Aktywne",
    value: "active",
  },
  {
    label: "Nieaktywne",
    value: "inactive",
  },
];

export default function Filters(): JSX.Element {
  const [selectedFilter, setSelectedFilter] = useState<string[]>(["all"]);

  const dialogRef = useRef<HTMLDialogElement>(null);

  const openDialog = () => dialogRef.current?.showModal();
  const closeDialog = () => dialogRef.current?.close();

  const toggleFilter = (filter: string) => {
    if (selectedFilter.includes(filter)) setSelectedFilter(selectedFilter.filter((f) => f !== filter));
    else setSelectedFilter([...selectedFilter, filter]);
  };

  return (
    <>
      <Button startIcon="mage:filter" onClick={openDialog}>
        Filtry
      </Button>

      <dialog ref={dialogRef} class={Styles["modal"]}>
        <div class={Styles["modal__grab"]} />

        <div class={Styles["modal__header"]}>
          <div>
            <h2 class={Styles["modal__title"]}>Filtry punktów</h2>
            <p class={Styles["modal__description"]}>Zawęź listę do tego, czego szukasz.</p>
          </div>

          <IconButton icon="material-symbols:close" onClick={closeDialog} />
        </div>

        <div class={Styles["modal_content"]}>
          {FILTER_OPTIONS.map((option) => { 
            const btnClass = selectedFilter.includes(option.value) ? Styles["modal__btn--active"] : "";

            return (
              <Button key={option.value} class={btnClass} variant={selectedFilter.includes(option.value) ? "contained" : "outlined"} onClick={() => toggleFilter(option.value)}>
                {option.label}
              </Button>
            );  
          })}
        </div>

        <Button class={Styles["modal_confirm-btn"]} fullWidth variant="contained" onClick={closeDialog}>Pokaż wyniki</Button>
      </dialog>
    </>
  );
}