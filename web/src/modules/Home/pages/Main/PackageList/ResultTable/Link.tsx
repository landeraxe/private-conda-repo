import React from "react";
import { Link } from "react-router-dom";
import styles from "./styles.less";

type Props = {
  channel: string;
  name: string;
};

export default ({ name, channel }: Props) => (
  <span>
    {channel} /{" "}
    <Link className={styles.link} to={`/p/${channel}/${name}`}>
      {name}
    </Link>
  </span>
);
