# DirHash
DirHash is a command line tool to take a directory and return the file hashes.

## Installation
At the time of writing, an installer is only available for amd64 Linux. Download the installer from the [releases page](https://github.com/melatonein5/DirHash/releases/tag/latest), and run `sudo ./dirhash_amd64_linux_installer'`.

Alternatively, you can download the binary and use DirHash the following way:

#### Running through the pre-compiled binary
At the time of writing, the binary is compiled for linux systems. Run DirHash file the following `./dirhash`.
```
Usage: ./dirhash [options]
Options:
  -i, --input-dir <dir>    Specify the input directory (default: current directory)
  -o, --output <file>      Specify the output file (default: no output file)
  -a, --algorithm <alg>    Specify the hash algorithm (default: md5), can take more than 1 argument, separated by spaces
  -t, --terminal           Output to terminal (default: false)
  -h, --help               Show this help message and exit
Supported algorithms:
  md5, sha1, sha256, sha512
Examples:
  ./dirhash -i /path/to/dir -o output.txt -a sha256
  ./dirhash --input-dir /path/to/dir --output output.txt --algorithm sha512 sha1
  ./dirhash -t
  ./dirhash --help
```

## Usage
Download DirFetch through the following command, and once DirHash is downloaded, move to the target directory using:
```
git clone https://github.com/melatonein5/DirHash.git
cd DirHash
```

#### Running through Go Runtime Environment (requires Go installation)
If a binary for your system is not available, you can run DirHash through the Go runtime environment with `go run dirhash.go`.
```
Usage: go run dirhash.go [options]
Options:
  -i, --input-dir <dir>    Specify the input directory (default: current directory)
  -o, --output <file>      Specify the output file (default: no output file)
  -a, --algorithm <alg>    Specify the hash algorithms (default: md5), can take more than 1 argument, separated by spaces
  -t, --terminal           Output to terminal (default: false)
  -h, --help               Show this help message and exit
Supported algorithms:
  md5, sha1, sha256, sha512
Examples:
  go run dirhash.go -i /path/to/dir -o output.txt -a sha256
  go run dirhash.go --input-dir /path/to/dir --output output.txt --algorithm sha512 sha1
  go run dirhash.go -t
  go run dirhash.go --help
```

#### Build from Source (requires Go installation)
To build from source, you can compile dirhash using `go build dirhash.go`, to compile an executable. The Go compiler selects your system architecture and operating system automatically. Note that if you have compiled to Windows, the `./dirhash` command becomes `dirhash.exe`.

## Example Outputs
There are two output types: terminal output; file output. Terminal output will show the results in formatted columns. File output will output to a file in csv format, regardless of file extension.

##### Terminal Output
```
[jmn@two DirHash]$ go run dirhash.go -a sha1 sha256
File Name                              Path                                                   Hash                                                             Hash Type
COMMIT_EDITMSG                         .git/COMMIT_EDITMSG                                    ac1f40bb9f2231060330d9c7a4ec9ab2355e300f                         sha1
FETCH_HEAD                             .git/FETCH_HEAD                                        5d037ee061137d7c6088c8a8b1f08b88d4899843                         sha1
HEAD                                   .git/HEAD                                              2aa05cb189709905d22504077e79b9d7ed74722a                         sha1
ORIG_HEAD                              .git/ORIG_HEAD                                         039435964f90dfd504265f1feb58e0eaff896bd2                         sha1
config                                 .git/config                                            5312c39011db3903bc8595f37f9ede2b761ec7b0                         sha1
description                            .git/description                                       9eca422d6263200fdf78796cab2765eb3fdc37e5                         sha1
applypatch-msg.sample                  .git/hooks/applypatch-msg.sample                       285716cd0b1f4e75e00d9854480fce108d1b6654                         sha1
commit-msg.sample                      .git/hooks/commit-msg.sample                           aa71dc4a856cae39c4da91d26e4cb1c7f0a8a92f                         sha1
fsmonitor-watchman.sample              .git/hooks/fsmonitor-watchman.sample                   c55bcc85dbbd3785080ff7d236a0102c47dfc5ba                         sha1
post-update.sample                     .git/hooks/post-update.sample                          f2ea063306cc8441634367a514d8fef76da8f45e                         sha1
main                                   .git/refs/heads/main                                   5d342e4e576ac445f7afc2e6c49a3e1986557019a5b1cb4acab447e5f3dbb73d sha256
HEAD                                   .git/refs/remotes/origin/HEAD                          6832d31adc6cd2a5714637a4395ad820983e5482535931d603d0261cae88b837 sha256
main                                   .git/refs/remotes/origin/main                          1378d3f729fbd8f5ec3d570cc94343c069ce13d0cbfb4d237e50f3f550980a68 sha256
LICENSE                                LICENSE                                                51b84b197025b72457944418eac94b71233e080bd0c45245198c7acad0e26aa9 sha256
README.md                              README.md                                              92f61da76c367c58263ace4d954095f6c86239d5ed0efb633ba203f082a689ae sha256
dirhash                                dirhash                                                eb567f761ed92ddb6a990498352bae2e8e878b99e638cc756e3a5ef5a7670a38 sha256
dirhash.go                             dirhash.go                                             8f2096ddeff9543ab56325dd63e9d689e61cd4b48ce074ea663c5a48cf58f5cd sha256
go.mod                                 go.mod                                                 40719d6dbfafc34704d3efd7c4bcea644755b48be7b1589e38cf86a4c7ffdf3b sha256
args.go                                src/args/args.go                                       86da89f3d83ed7fe84a62c946c72401df10274d73c5543a0590cfa0569b1c38e sha256
hash_algorithm.go                      src/args/hash_algorithm.go                             64b723c3cc795a25397cd529cc9148e8b53bc5de83d4750a31fedafe7e9a2d05 sha256
parse_args.go                          src/args/parse_args.go                                 92f5a1bd99e8659cdca3689958861018425611d72691471a050b88385da22281 sha256
help.go                                src/cmdline/help.go                                    b17700c422f02e6de8b2f7f6bac4208600efba96c48e444af7e38c440c6ece10 sha256
output_files.go                        src/cmdline/output_files.go                            e0fc4f16df0961806b235c17bb363f33ab985be7b5e37a48075e95d56dc0186a sha256
enumerate_files.go                     src/files/enumerate_files.go                           c121e905cb787652d97423af50e5f68b67fd62b5d0ea93315a0ca79334830e58 sha256
file.go                                src/files/file.go                                      b60cff7dcd0a219b435ef39e5905796885f51ed3bc0febc0e41befab336374aa sha256
hash_files.go                          src/files/hash_files.go                                30f098d18815d96bc3a10b114486cebcab0e8725e98681d22baced6a508ff752 sha256
md5_files.go                           src/files/md5_files.go                                 43f7e1c3b48076af3ffc6ca10572528e697db175cdf24e1533bf73ba0a1241b0 sha256
sha1_files.go                          src/files/sha1_files.go                                8c54814b5b021ae6085ae91f7bf2c923dd8fe9ff590ee13946f7086c14d5f87c sha256
sha256_files.go                        src/files/sha256_files.go                              2e14eea5233669ffa6c32e8d91e5ce1baaef6aa4cac48bda82607264c77bbf07 sha256
sha512_files.go                        src/files/sha512_files.go                              bb427b881944e600d521532200fb5106b240f940b1d899067ff30832ab667abb sha256
write_output.go                        src/files/write_output.go                              976f2684e4b610428bce50ece072c76351f4cf476f237825ebb8f01b02310a26 sha256
```

#### CSV Output
```
Path,Hash,HashType
.git/COMMIT_EDITMSG,34bb1ef88028a1b61d48100cc9e1dfa37b8af19f1e3babe2747d4f03dbb9d4ac4f464cb01e8e11817e6b2050661a6d3f6c8da21bd5608e236bfcd3f87083b9a7,sha512
.git/HEAD,d2af064a992de738fe2445adddeb0fb6085e74a698937dacaed294ae29376c869bdaa9c0ca1453f8014c168368916dcdf851a595b8515515d7f8f750fbd2acdb,sha512
.git/config,686dee4a9d02132ad1d5f50843952685a330d8e9f6b65fd21cd3828d690695f61f88597cb4b913ec8d8673e5349dd7fefa6d91c2ae8939c28449810bab7e5747,sha512
.git/description,6b7a652443e815dee8d16e71ac80625fe78affeaf54ba190ac229950ba9072bd89c90a5d2a01224070c827358ff185f57a2efe7abd2b1333cb5da84102c2e784,sha512
.git/hooks/applypatch-msg.sample,9f2d3ac648424b9c6a2e462837ef0666d4d89820921fec173a55a003c941b63b2601cae0eb57938b2056ccb571d5df8bbfd2617bb3a8bd1c8f0c8607beb2b837,sha512
.git/hooks/commit-msg.sample,8d5395d45e8718a94ffa766b9cdbf064ef31557fdbd22cffb090fb9ccf41268b837b5efa84742272bf2c95e97fd1a7146a6a3cf00de86caf09440fa63bae84ec,sha512
.git/hooks/fsmonitor-watchman.sample,c1ace3bc31e0fe3076f0b97d7fe89c416f339b894f925539a7c358d697ddc3e3d2e25740c8c2da41018f16e419f954e67adbd4bc2129ecef5b798b951a6d0a62,sha512
.git/hooks/post-update.sample,73ca6bfcb2e478fe87cc2f1e7dff55bb3f2588d7785b3b4bc180c83b75ad4cb456c6d3fdd370b52a1a4887f250db51ea7ceac5fad12a8e519f3b5ebcfc77a6f4,sha512
.git/hooks/pre-applypatch.sample,22e9e24f8343baef07bcfe7e412b3c67755f3d3eb5445193419340091f72575722e4a906a5aa2115f44aa169e64fbf5c237f6ab5cf6fba054af4a7252d7ae5d2,sha512
.git/hooks/pre-commit.sample,8587e03cd0a452bc917efc34af72c2876be7d4adf7347cdc4db37338ee283b0d3b04b4a111be430bc851e0b27888e2e96a00807a368d77730e9db5c1a4af166c,sha512
.git/hooks/pre-merge-commit.sample,8422a1c556ba91a12848471aeb928ed51c3d77f03e3e3b71cd637b40dca25ca4f23c4ffe51e8d099dba9b09c368aff0577fe7e23feb3e1fccab1cef27029a39b,sha512
.git/hooks/pre-push.sample,8ad5fdf9fd482f5909e8375fc647f0a91e1b15de83e427609127632b9b9b892c3d5c4bacf8280ebba1c94d589fc38fead3d74145862514a8326223da70ef51b7,sha512
.git/hooks/pre-rebase.sample,11ed8d5a3ebdfea1ec098aae2a39dda4bf673f3e04a841a798ed4dbc76f628ae293b2bab8ec3a493aa28e4d822fa8e9d56fafc9aad8b285d26d2923b33ce827d,sha512
.git/hooks/pre-receive.sample,6ff196f0f66524e96cb11200fdee957051de399d0519c3ec1d1d4826bdd09fe4a9bf81ac6201c6c266bf30eea2c1d069ae2cb4c504e0dad9aee0055852243a8b,sha512
.git/hooks/prepare-commit-msg.sample,a3e4495f2e707e714f3e3342b3d90a695b1f579ae2f442c057f415f18b25935a67a46e9503ea007476e23eaf3b1a4a7a8d7b725aadd4d0b8d613ff413521ea0d,sha512
.git/hooks/push-to-checkout.sample,207eb72a68379db76dc6a01440448bfc4b3c5ac5f1ab31ab5179bad9bc6adb066a275e892dc6d5cca020564345f512efdea11d008120ef5853185ca1beecf526,sha512
.git/hooks/sendemail-validate.sample,a858731c5fc2e87a7155d4f9f4be0730af763e7d8e71fac154f5303e2ded716308728dd2583449461a40688eeff618cf62dd388ff61ac61ce9d7daebc9ca8698,sha512
.git/hooks/update.sample,a1910b35e0af53a3b5bba0872ae2198d55962f2611692f0dc0ada15f6789f674655673810c065c869f86be723c2b0a33a6efc7d5ebeac2cd1abcdbc99badc737,sha512
.git/index,f3a27d665ec7d39db84e7b7bf9621269240a0609c6c7d7987dfdc870f83b2c6d786ec1781bf0bb4d73107c94ba1144d311bb875bdb09629da1d609d132b48675,sha512
.git/info/exclude,0f018163e053ef1c6e75f00fee2ac08e38795104722574d85901a3e05d73d1df17714a6053e834b8936d2dd69c7b435e171b5ccffbfa04d39c04810ac1522208,sha512
.git/logs/HEAD,f6d2b6cbf9e4f5c8440c2b4934b25f997807d181bd09fc18d9db3b9a65e827ff2effe09bc3940132c1f470d5129054ecb8833bb21b800c9e80fe8468f6ba380b,sha512
.git/logs/refs/heads/main,11fbbc6c622a1618edf5a8814fbf7707fab2174d7b3c34f2d311df06877ebea684858f9204e52889b3282359ee6f955f4859295d9cdadb925723debd045f8489,sha512
.git/logs/refs/remotes/origin/main,ed8a3510ee878b1e37793b9a1386e61e97fcf583895bcb1f96b74ac662a67a3619cb3a438a38671124278d8556f21e473adf4331afe9167568105ccfad7cca77,sha512
.git/objects/2f/ffd99de246994a595d09f719540338bc102824,fca87af0078af2d48e7c42d423af6c191ad7180a337a99d8b493cd2ba3dfcb3bb8cf2b7be8329580d53a3392b412a71a99905b953183a348387a8fbe3cffdd7f,sha512
.git/objects/79/7f0061b697c20dc65f0abe1be1b1881d222c75,9892bd79feb79d1e5f051aaba40a7c47936e88ce400bce14d7708ce2016bff8b9decd00988be62a7509cfc3d9507df09fbbc099b3c7664bed8982558bd6be795,sha512
.git/objects/bb/3b39e44f6d27e2f6367953794fc03abb1b89cb,1a9a7b4ab724b3a2937213136a22e4f03d9150e72231749df1c1a53cf62fdcd130727c066c9173f6e81efd5e3b2de0eb7defdc688affd3f12100a98753bd23aa,sha512
.git/refs/heads/main,fc935d2217d4f30f4d55bde62503063330b5846164c7c8f1934355e9173b2a585cce2d094e9442577fb0600dea53dc2e458d547c70889e088238c974c4823134,sha512
.git/refs/remotes/origin/main,82c00459623eb0939e36ce8b34f7a15b4cf14da47bc561d65f6757178fad9fd5a300da29209da99cf0993d3bc9446ffa69ab7fefd92ac77eb1afba0e955141bc,sha512
README.md,827d8ea799490cb1d8721239a9dd719e0d6de3afb661a724f601e8ee4f99b399a2216da83e432a5b16aae6abb3b88aa8c32a5226a49b457b098c248f3e157588,sha512
dirhash,9f63843a3082cbaea927b162e40241a7a91990fd25ddb6ffb6667964abe82a09b7d7bfdcb0ff813b712cf5b29fb93c305924badf9b06d52a9093673290f1b1d2,sha512
dirhash.go,6f5339d4f462a0a363d58c05afe747937deb5cffbb9839e6a2cc9431178f104e45ca27118101a34fe086e0a7001fb58d155251bf52d6bd4c27472544f2b94cf3,sha512
go.mod,ad6773ab853200b272eec31380185a41d9f832878fd241764756b9a68cb71f971d1a80271fb9ac3a9f596e3ae9620c87c995fe84423eec8e64afebe2a0ad2f7d,sha512
src/args/args.go,42772e9f5c9730dad50ef04bb086d2f29d8974bea9946d51bd1401fe39a15151ae5a1a98778542d35b4ff9fd83f856fb86ddaa5ac19f7d683a0df015e5b0f793,sha512
src/args/hash_algorithm.go,519049bd042054c3a1158d27fc259fa1d53feac68c8a1db837c8f9049dd4257d5d08cdba997bfc05ec804aa01e221538d04db4c7c4e431aba19c4eb5666af2ab,sha512
src/args/parse_args.go,5245d2ec0c3ec77fc97d9a89b1a1a6b10df64aea82646790b5cff925bbb9ca8bde8ed3df7c9c505727ab0f0c36479c30f4db081fbd0acaccd40f691681c596ac,sha512
src/cmdline/help.go,9a0d3d9bd9e73dc4855b78f8ff8dd7db1c623a07563b404a976071097a6e4ba39c39d1e08c526fec1f009b836808951555c37ff90f83237b14bff56601427acc,sha512
src/cmdline/output_files.go,3dd2c94e2b90586973d8d65169efbceeb730b28faae7510bfd63d0f8a24cab76cefc27279d9a428e9c797939f40f723f63f416ee51002c03105ce4afe341efe2,sha512
src/files/enumerate_files.go,232aa8c8a49ce73db529baede71bb9e3c401693a481171902f6406d27f3fcfe1a260ac9316802830f94cb308c4373237bd56d718c9fa3c9f8ce59d69b5036529,sha512
src/files/file.go,670ba48849464991f8070ebfcffd2bb5a86781faae38492cab56331be9803a549a9064ee4def4865872bb545eb7a2ace4d2d94e8f3495c9f9db765d62a6f3d4e,sha512
src/files/hash_files.go,6865f9a842998a822dd4ed481c86a2d0081e37998d0bf1681ed0133f333d916330cd7a47a6eaf1c78dbddcdd083164125191be240c620da2e260ac0a82b19416,sha512
src/files/md5_files.go,e03534dd45d373e743a65bf622bd3815286f101092c713f17df53baf4341259165969ee4f9409e5dd2fe64cd55660dd926fe069bb34bfa4442fff77a5555b729,sha512
src/files/sha1_files.go,e66eb9f2a87e21051f46b6590f6e7cd5572fa35b74024a35728a36f8c61d2421b9155f3cbec67fcc3e87bed26029f6418b218aeefa022446af35b1727eb3f32f,sha512
src/files/sha256_files.go,182f1a84ad40d3c4dc50aa9dc81e7e316aae6172a7a66621052765d296ba0ae889fd4c3c59e60d784754f2cdd39d6bd320fa52e420c98ce0a2816fc87cd36b0a,sha512
src/files/sha512_files.go,b03f1a3d9e3f64f11f81b4011ff3579f7b8a93ea97c87ac1b5544fe9e5e03b054cbc37650aadec7b93e31f2e29ceb750b89fd268c5a424188bb1a27235bcab08,sha512
src/files/write_output.go,3eb5bb989a47bfc677027578be94f8d3190014c182839aa5eb7cb08a6ad9ddb3c29926d95c9616e1a58ebc9accacf9a3463c36ea844dd22342437b053c207008,sha512
```

## Roadmap
The main improvements to make are: 
1. Support for more hashing functions.
