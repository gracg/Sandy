package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func ClientMain(cnfg Config) {

	for {
		client := &http.Client{}

		req, _ := http.NewRequest("GET", cnfg.Client.Target, nil)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		if err == nil {
			if resp.StatusCode == http.StatusOK {
				respBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}

				var keys ResultKeys

				err2 := json.Unmarshal(respBytes, &keys)
				if err2 != nil {
					panic(err2)
				}

				var approvedKeys []string
				for x := range keys.Keys {
					z := keys.Keys[x]
					_, err = VerifySignedKey(z.Key, z.Signature)
					if err == nil {
						approvedKeys = append(approvedKeys, z.Key)
					}
				}

				if len(approvedKeys) > 0 {
					AuthKeyFileExists, _ := exists(cnfg.Client.AuthorizedKeyFile)
					if AuthKeyFileExists == false {
						f, err := os.Create(cnfg.Client.AuthorizedKeyFile)
						if err == nil {
							f.Close()
							os.Chmod(cnfg.Client.AuthorizedKeyFile, 0644)
							ioutil.WriteFile(cnfg.Client.AuthorizedKeyFile, []byte(parseKeys(approvedKeys)), 0644)

						}
					} else {
						ioutil.WriteFile(cnfg.Client.AuthorizedKeyFile, []byte(parseKeys(approvedKeys)), 0644)

					}
				}
			}

		}

		if resp.StatusCode == http.StatusOK {
			time.Sleep(time.Second * 30)
		} else {
			time.Sleep(time.Second * 5)
		}
	}
}

func parseKeys(k []string) string {
	var s string
	for x := range k {
		z := k[x]
		z = strings.TrimSuffix(z, "\n")
		s += z
		s += "\n"
	}
	return s
}
