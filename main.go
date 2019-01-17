package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const WeekToUnix = 604800

type ResponseListFriends struct {
	ResponseListFriends ResponseListFriend `json:"response"`
}

type ResponseListFriend struct {
	Count int32  `json:"count"`
	Items []Item `json:"items"`
}

type Item struct {
	Id        int32    `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	LastSeen  LastSeen `json:"last_seen"`
}

type LastSeen struct {
	Time int64 `json:"time"`
}

func main() {

	StartStrippingFriends()

}

func StartStrippingFriends() {

	accessToken, numberWeeksUserNotLoggedIn := ReadInputData()

	searchQueryForAllFriends := CreateSearchQueryForAllFriends(accessToken)

	respBodyByte := ExecuteCreatedRequest(searchQueryForAllFriends)

	allFriends := ParseResponseJSON(respBodyByte)

	suitableUsers := CheckUsersByParameters(numberWeeksUserNotLoggedIn, allFriends)

	DeleteUsersFromFriends(suitableUsers, accessToken)

}

func ReadInputData() (accessToken string, numberWeeksUserNotLoggedIn int64) {

	fmt.Println("Вставьте свой access_token:")
	_, err := fmt.Scan(&accessToken)
	CheckErrors("func StartReplaceFormatNEL()", err)

	fmt.Println("Введите кол-во недель, которое не заходил пользователь в сеть:")
	_, err = fmt.Scan(&numberWeeksUserNotLoggedIn)
	CheckErrors("func StartReplaceFormatNEL()", err)

	//TODO: Add Check input data

	return
}

func CreateSearchQueryForAllFriends(accessToken string) (urlRequests string) {

	beginningUrlRequests := `https://api.vk.com/method/friends.get?count=5000&order=name&fields=last_seen&access_token=`
	endUrlRequests := `&v=5.92`

	urlRequests = strings.Join([]string{beginningUrlRequests, accessToken, endUrlRequests}, "")

	return
}

func ExecuteCreatedRequest(myUrl string) (respBodyByte []byte) {
	resp, err := http.Get(myUrl)
	CheckErrors("func execGET(myUrl string) http.Get(myUrl)", err)

	defer func() {
		err = resp.Body.Close()
		CheckErrors("func execGET(myUrl string) resp.Body.Close()", err)
	}()

	if resp.Body != nil {
		respBodyByte, _ = ioutil.ReadAll(resp.Body)
	}

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBodyByte))

	CheckErrors("func execGET(myUrl string) resp.Body.Read(bs)", err)

	return
}

func ParseResponseJSON(respBodyByte []byte) (allFriends []Item) {
	var responseListFriends ResponseListFriends

	err := json.Unmarshal(respBodyByte[:], &responseListFriends)
	CheckErrors("func ParseResponseJSON(respBodyByte []byte) (allFriends []Item), json.Unmarshal(bs, &responseListFriends)", err)

	allFriends = responseListFriends.ResponseListFriends.Items

	return
}

func CheckUsersByParameters(numberWeeksUserNotLoggedIn int64, allUsers []Item) (suitableUsers []Item) {

	inputWeekUnix := numberWeeksUserNotLoggedIn * WeekToUnix
	unixTimeNow := time.Now().Unix()

	for i := range allUsers {
		if unixTimeNow-inputWeekUnix > allUsers[i].LastSeen.Time {
			suitableUsers = append(suitableUsers, allUsers[i])

			fmt.Println("Имя, Фамилие:", allUsers[i].FirstName, allUsers[i].LastName)
			fmt.Println("Id пользователя:", allUsers[i].Id)

			fmt.Println("Последний раз был в сети:", time.Unix(allUsers[i].LastSeen.Time, 0))

			currentUserId := fmt.Sprint(allUsers[i].Id)

			linkToUser := strings.Join([]string{`https://vk.com/id`, currentUserId}, "")

			fmt.Println("Посмотреть пользователя по ссылке", linkToUser)
			fmt.Println("------------------------------------")
		}

	}
	fmt.Println("Кол-во пользователей рекомендуемых к удалению:", len(suitableUsers))
	return
}

func DeleteUsersFromFriends(friendsToDelete []Item, accessToken string) {

	userIdUrlRequests := `https://api.vk.com/method/friends.delete?user_id=`
	accessTokenUrlRequests := `&access_token=`
	endUrlRequests := `&v=5.92`

	for i := range friendsToDelete {
		urlRequests := strings.Join([]string{userIdUrlRequests, strconv.Itoa(int(friendsToDelete[i].Id)), accessTokenUrlRequests, accessToken, endUrlRequests}, "")

		fmt.Println("urlRequests", urlRequests)

		resp, err := http.Get(urlRequests)
		CheckErrors("func DeleteUsersFromFriends(friendsToDelete []Item, accessToken string) http.Get(urlRequests)", err)

		var respBodyByte []byte

		if resp.Body != nil {
			respBodyByte, _ = ioutil.ReadAll(resp.Body)
		}

		resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBodyByte))

		fmt.Println("respBodyByte", string(respBodyByte))
		err = resp.Body.Close()
		CheckErrors("func execGET(myUrl string) resp.Body.Close()", err)

		currentUserId := fmt.Sprint(friendsToDelete[i].Id)
		linkToUser := strings.Join([]string{`https://vk.com/id`, currentUserId}, "")

		fmt.Println("Успешно удалён:", linkToUser)
	}
}

//Общая проверка всех ошибок
func CheckErrors(methodName string, err error) {
	if err != nil {
		log.Println(methodName, "get errors:", err)
	}
}
