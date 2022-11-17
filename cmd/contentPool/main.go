package main

import (
	"fmt"
	"git.in.zhihu.com/zhsearch/search-ingress/pkg/common/ztext2"
	"regexp"
	"strings"
)

const Template = `\<[\S\s]+?\>`

func main() {
	//	buf := `1.其实雨不大,是风搞得紧张。
	//2.不需要建议,不需要认同,喜欢就值得。
	//3.求之不得往往不求而得。
	//4.五颜六色的生活,不能乱七八糟的过。
	//5.风说了许多,把夏天注的盈满。
	//01.考虑一千次，不如去做一次；犹豫一万次，不如实践一次；迈出第一步，才有可能成功。
	//02.真正改变命运的，并不是我们的机遇，而是我们的态度。
	//123.,你说的爱太简单了 我要的是公螳螂被吃掉的心甘情愿 是到死都牵挂 是生生世世的宿命不息
	//124.,你老爱说自己不好 说自己见过好多人 那天我无意间侧头看你的时候 你眼神流露的分明是想好好被爱
	//`
	//s := "<p data-pid=\"ncPgVsFL\">1.</p><blockquote data-pid=\"tVGkw1gJ\">一定要收到来自夏天的惊喜 那就是录取通知书</blockquote><p data-pid=\"WgMZyojF\">2.</p><blockquote data-pid=\"LeCI0T10\">改天是哪天，下次是哪次，以后是多久，去经历，去后悔。保持热爱，奔赴山海。</blockquote><p data-pid=\"rszjjn0G\">3.</p><blockquote data-pid=\"1gErxRzj\">你一定要走，走到灯火通明。</blockquote><p data-pid=\"R4CnJfxW\">4.</p><blockquote data-pid=\"vH3sXE8z\">夏天是冒泡的冰可乐，酸奶味的棒棒冰和捧在手心的冻果冻，是简单的白衬衫和放在球场边的矿泉水，是温柔的情歌和听不够的rap</blockquote><p data-pid=\"_tvs-oYq\">5.</p><blockquote data-pid=\"zFAAHgTL\">温柔的日落中总要夹杂些诗和远方</blockquote><img src=\"v2-529ab90d5550696a3a3fbfd0bc7e4397.jpg\" data-caption=\"\" data-size=\"small\" data-rawwidth=\"690\" data-rawheight=\"519\" data-watermark=\"original\" data-original-src=\"v2-529ab90d5550696a3a3fbfd0bc7e4397\" data-watermark-src=\"v2-6593e7add08366df6d07aee9f4df66b3\" data-private-watermark-src=\"v2-41d08bed22f39f8b85a6aca86ccce530\"><p data-pid=\"BTjSybFh\">6.</p><blockquote data-pid=\"2Yi7yUcm\">你好我是快乐警察请您配合我的工作我现在要抓走您所有的不开心</blockquote><p data-pid=\"2Glq4Pk1\">7.</p><blockquote data-pid=\"R6mo2BAM\">你也一定会前程似锦</blockquote><p data-pid=\"pKd8eMSZ\">8.</p><blockquote data-pid=\"_7ZcmvZr\">黑马诞生 金榜题名 祝你也祝我</blockquote><p data-pid=\"G80z6ePI\">9.</p><blockquote data-pid=\"Az1SvozH\">说到做到这个词汇既踏实又浪漫</blockquote><p data-pid=\"2sHwAlqD\">10.</p><blockquote data-pid=\"XBQZRwAf\">像我这么单纯的人做不来这么有心机的数学题</blockquote><img src=\"v2-e77fbc2ea245854a1a9539a373931b06.jpg\" data-caption=\"\" data-size=\"small\" data-rawwidth=\"690\" data-rawheight=\"519\" data-watermark=\"original\" data-original-src=\"v2-e77fbc2ea245854a1a9539a373931b06\" data-watermark-src=\"v2-9a673cbe50ea8b194d47b7e1eebff7bb\" data-private-watermark-src=\"v2-2100807812726d2430909d6bcf0ec485\"><p data-pid=\"HSDwT-zK\">11.</p><blockquote data-pid=\"arcYirjx\">一定要站在你所热爱的世界里闪闪发光。</blockquote><p data-pid=\"F-Hd2nYH\">12.</p><blockquote data-pid=\"LUVhJvur\">听说神不能无处不在，所以创造了妈妈。到了妈妈的年龄，妈妈仍然是妈妈的守护神。</blockquote><p data-pid=\"rj1gmGhT\">13.</p><blockquote data-pid=\"KzRviDU8\">这个世界太危险，时间就该浪费在美好的事物上。</blockquote><p data-pid=\"-GQIviL1\">14.</p><blockquote data-pid=\"2uJ6bxIE\">爱自己是终生浪漫的开始。</blockquote><p data-pid=\"Eev_qv2m\">15.</p><blockquote data-pid=\"7uiqbthL\">煎和熬都是变美味的方法，加油也是。</blockquote><img src=\"v2-df0487bf0a2604678ffe50925f33b3f1.jpg\" data-caption=\"\" data-size=\"small\" data-rawwidth=\"690\" data-rawheight=\"519\" data-watermark=\"original\" data-original-src=\"v2-df0487bf0a2604678ffe50925f33b3f1\" data-watermark-src=\"v2-50659c80d807f1639a79511a2dd08ce9\" data-private-watermark-src=\"v2-a20dcf3ae5e95667fd922a22cdd254c2\"><p data-pid=\"lxEK2weA\">16.</p><blockquote data-pid=\"WYq_hx0d\">夏天的避难所当然是图书馆啦。</blockquote><p data-pid=\"XtVEakdo\">17.</p><blockquote data-pid=\"CbjoIZ81\">世间风物论自由，喜一生我有，共四海丰收。</blockquote><p data-pid=\"YOpGtf5B\">18.</p><blockquote data-pid=\"kWQrOGh8\">原来躲起来的星星也在努力发光啊。</blockquote><p data-pid=\"u2mV4bwo\">19.</p><blockquote data-pid=\"JecC-Vha\">你是要长大的小朋友， 一定要撑着， 别倒下去。</blockquote><p data-pid=\"LLbQo8eb\">20.</p><blockquote data-pid=\"JIYqQy57\">只要我能拥抱世界，那拥抱得笨拙又有什么关系。</blockquote><p data-pid=\"0UeWydGu\"><b>更多宝藏文案还请关注公众号：木岁说 一个日更的微信公众号</b></p><img src=\"v2-416bbb9ef98b719cb89e37365b4bbdad\" data-caption=\"\" data-size=\"normal\" data-rawwidth=\"1125\" data-rawheight=\"529\" data-watermark=\"\" data-original-src=\"\" data-watermark-src=\"\" data-private-watermark-src=\"\"><p></p>"
	s := "<p data-pid=\"P17QQFkT\">1.普吉岛的夏天永不停歇，我热爱的少年永远热恋！</p><p data-pid=\"KHJCuNcZ\">2.少年的心动是仲夏夜的荒原，割不完，烧不尽，长风一吹，野草就连了天。</p><p data-pid=\"kStvryX8\">3.那年的盛夏比往常来的聒噪，教室外枝桠疯长，却总挡不住烈阳。</p><p data-pid=\"s74xfhU3\">4.那是系统里永远看不到的景色，是万家灯火，是喧嚣人间。</p><p data-pid=\"zWgWZLZo\">5.一星陨落，黯淡不了星空灿烂 ；一花凋零，荒芜不了整个春天。</p><p data-pid=\"0Wc7myfU\">6.这个世界乱糟糟的 而他干干净净 可以悬在我的心上 做太阳和月亮</p><p data-pid=\"FufmYr8M\">7.我曾在人山人海的星光里为你撕声呐喊，最美好的岁月都和你有关</p><p data-pid=\"kw7MP2xx\">8.“他踏光而来，我怎配让他落入尘埃.”</p><p data-pid=\"ra9mlgaE\">9.“再红一点吧 红到路人皆知 红到没人敢诋毁 红到大街小巷处处都是粉丝 红到被世界认可&#39;&#39;</p><p data-pid=\"1kJzkNNf\">10.漫漫追星路，幸遇意中人。</p><p data-pid=\"weKg4VOy\">11.从遇见你开始，所有的晦暗都留给过往，凛冬散尽，星河长明。</p><p data-pid=\"SyyLeeGh\">12.人间美好大概就是：傍晚的日落和柔情的你</p><p data-pid=\"9Vk8dWxR\">13.“在追逐月亮的途中 我也被月光照亮过”</p><p data-pid=\"BR83ONFr\">14.瞒着浩瀚宇宙向你靠近，星河万顷是你的见面礼</p><p data-pid=\"mYqaB7-y\">15.我将告别黄昏，挣脱藏身的黑暗，想你的光里坠落</p><p data-pid=\"Lhlda-4w\">16.浩瀚宇宙无法私有，但却可以寄存于追求理想的眼中。</p><p data-pid=\"yzyZIUqv\">17.所有的苦难和背负尽头都是行云流水般的此时光阴。</p><p data-pid=\"hRpxYDVA\">18.走过黯然无光的黑夜，才能看见黎明</p><p data-pid=\"l69BXYH-\">19.你在黄昏寄一场梦给我，满载相思的欢喜和宇宙的星辰</p><p data-pid=\"fhiDUU4w\">20.少女的征途是星辰大海而并非烟尘人间</p>"
	//带空格
	v := "1 “夏天好像总是热烈和明亮 天空晴朗 月亮皎洁 请你也要万事顺意 奔赴更好的地方 在这个盛夏。” 2 “相逢”是个很美好的词儿，因为它总是会带给人不期而遇的温暖和生生不息的希望。 ” 3 “如果我写不出温柔漂亮的句子，你能不能直接看一看我温柔的心。” 4 第一眼就心动的人 再看多少眼都还是会心动 5 “只有频率相同的人 才能看到彼此内心深处 不为人知的温柔 理解你的山河万里 尊重你的不同选择 ” 6 “让我感觉最美好的是下课和姐妹一起冲去食堂 晚饭后手拿牛奶一起逛操场聊喜欢的男孩 体育课坐在篮球场外的草坪上 偷偷看喜欢的人打篮球的时光” 7 “离睡着还有一万九千八百七十次想你。” 8 “我真的十分喜欢拥抱了 被一把抱住时 阳光 洗衣粉 所有温暖的味道 都好像一下子融进了身体里 倘若下次再见面 希望你二话不说过来抱我。” 9 “有人爱晚风中绽放的花 小巷里乱窜的猫 独自一人看的日落 真诚而不失浪漫地热爱着自己所热爱的” 10 “月亮月亮，你能照见南边，也能照见北边，照见他，你就跟他说一声，就说我想他了。” 11 “年少时的朋友没能一路走下去多少有点可惜，可这一生还长着呢，谁知道下一个遇见的会是怎样奇妙的一个人儿。” 12 “人最勇敢的是：哪怕知道自己会受伤，可是我们还会选择勇敢去爱。” 13 ᴹᵃʸ ʸᵒᵘʳ ˡⁱᶠᵉ ᵇᵉ ᵃˡʷᵃʸˢ ʷᵃʳᵐ ᵃⁿᵈ ᵍᵉⁿᵗˡᵉ ᵃⁿᵈ ˢʰⁱⁿⁱⁿᵍ. 愿你的生活常温暖，日子总是温柔又闪光。 14 ᴳᵒ ˢˡᵒʷˡʸ, ᵗʰᵉʳᵉ ʷⁱˡˡ ᵃˡʷᵃʸˢ ᵇᵉ ˢᵒᵐᵉᵒⁿᵉ 慢慢走吧 总会有一个人和你步伐一致的 你不用去追他 他也不用去等你 15 ʸᵒᵘ ʷᵃⁿᵗ ᵗᵒ ᵇᵉ ᵇᵉᵗᵗᵉʳ, ⁿᵒᵗ ᵃˡʷᵃʸˢ ˡᵒᵒᵏ ᵇᵃᶜᵏ! 你要去变得更好，而不是总回头！ 16 ᴵᶠ ʸᵒᵘʳ ᵏⁿⁱᵍʰᵗ ᵈᵒᵉˢ ⁿᵒᵗ ᶜᵒᵐᵉ, ʸᵒᵘ ʷⁱˡˡ ᵃˡʷᵃʸˢ ᵇᵉ ᵐʸ ᵖʳⁱⁿᶜᵉˢˢ. 你的骑士没有来，你就永远是我的公主。 17 “谁都无法给你未来，因为你才是自己的未来。” 18 人们从诗人的字句里，选取心爱的意义。但诗句的最终意义是指向你。 ——泰戈尔《吉檀迦利》 19 “不是只有礼物和鲜花才浪漫，愿意听我碎碎念也很浪漫。” 20 “有时候我看着夜空，想到我们至少看到的是同一轮月亮，就觉得我也并非错过了一切。” 21 “好好生活，怀念的不一定要见面，喜欢的不一定要在一起，每一种距离和遗憾都有它存在的意义。” 22 “我以为一个人要很优秀，才能被人喜欢。其实只要一个人足够热烈，就一定能点燃另一个人。” 23 “ 再努努力，你要相信山顶远比想象中更美。” 24 怎么解释你的温柔程度呢？你打了个呵欠，宇宙就老去一点。你多看我一眼，岁月便又，重来一遍。 ——李倦容 以上图cr微博：-初夏未满 25 我爱你，不是因为这里只有你一个人而爱你，而是因为只想爱你一个人 ——《情书》 26 “喜欢那些精致的女孩子 出门化妆很久自拍很久 看似自恋且浪费时间的行为 在我眼里则是充满认真生活的浪漫” 27 “你是一朵特别的花，想怎么长就怎么长，不一定非要长成玫瑰，就算全世界都在期盼你长成玫瑰，我还是希望你可以做自己。” 28 “人生总是奇妙的，一旦你努力去做一件事，如果结果不是你想象的那样，那么老天一定会给你一个更好的结果。” 29 以前总觉得万物皆有遗憾，但现在突然想开了，一切都是相互选择的结果。我尽力了遗憾的不该是我。 30 总有一天我会带你去海边吹吹晚风 看看日出日落 当这些只想与你一起才有意义 31 “请在五月，在夏天来临之前变得更厉害一点吧，要在夏天拥有很多的开心和快乐。 ” 32 希望一如既往的坚强 站在迎着光的地方 活成自己想要的模样 我的原创公众号～欢迎关注 碎碎念： 就到这里啦，好像好久好久没有发温柔小句子了。夏天好热，最近生活也不太安然，但浪漫至死不渝。\t"
	reg1 := regexp.MustCompile(`\d+[”、".． ]+(.*?)+\s`)
	reg2 := regexp.MustCompile(`\d+[”、".． ]+(.*?)+\d`)
	//reg1 := regexp.MustCompile(`^\d+[”、".“]+(.*?)+[”。“、".]+\s`)
	//reg2 := regexp.MustCompile(`\\<[\\S\\s]+?\\>`)
	if reg1 == nil {
		fmt.Println("regexp err")
		return
	}
	r := regexp.MustCompile(Template)
	fmt.Println("-------", ztext2.ZTextPlaintext(s), "\n", strings.Replace(ztext2.ZTextPlaintext(s), "\t", "", -1), "\n", r.ReplaceAllString(strings.Replace(ztext2.ZTextPlaintext(s), "\t", "", -1), "\n"))
	//result1 := reg1.FindAllStringSubmatch(r.ReplaceAllString(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(v, `“`, "", -1), `”`, "", -1), " ", ",", -1), "\n", "", -1), "\t", "", -1), "\n"), -1)
	result1 := reg1.FindAllStringSubmatch(strings.Replace(ztext2.ZTextPlaintext(s), "\t", "", -1), -1)
	result2 := reg1.FindAllStringSubmatch(r.ReplaceAllString(s, "\t"), -1)
	result3 := reg1.FindAllStringSubmatch(v, -1)

	result4 := reg2.FindStringSubmatch(v)
	//result5:=strings.Match(/d(\S*)/)[1]
	//result2 := reg2.ReplaceAllString(s, "")
	//fmt.Printf("%+v\n", result1)
	for _, strings := range result1 {
		fmt.Println("!!!!!", strings[1])
	}
	for _, strings := range result2 {
		fmt.Println("。。。。", strings[1])
	}
	for _, strings := range result3 {
		fmt.Println("333333333333", strings[1])
	}
	fmt.Printf("%+v\n", result4[len(result4)-1])
	fmt.Println("@@@@@@@@@@@222")
	a := strings.Trim("端午节祝福", "端福")
	fmt.Println(a)
	fmt.Println(strings.Trim(a, "祝福")) //去掉字符串s中首部以及尾部与字符串cutset中每个相匹配的字符
	fmt.Println(strings.Trim("Hello, Gophersh", "Hh"))

}
